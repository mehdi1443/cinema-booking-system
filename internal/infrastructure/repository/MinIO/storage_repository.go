package MinIO

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type S3Storage struct {
	env      *bootstrap.Environment
	client   *s3.S3
	uploader *s3manager.Uploader
	buckets  map[enum.BucketType]string
}

func NewS3Storage(env *bootstrap.Environment) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(env.Storage.AccessKey, env.Storage.SecretKey, ""),
		Region:           aws.String(env.Storage.Region),
		Endpoint:         aws.String(env.Storage.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	buckets := make(map[enum.BucketType]string)
	buckets[enum.MoviePoster] = env.Storage.Buckets.MoviePoster
	buckets[enum.MovieTrailer] = env.Storage.Buckets.MovieTrailer
	buckets[enum.ActorPhoto] = env.Storage.Buckets.ActorPhoto
	buckets[enum.CinemaPhoto] = env.Storage.Buckets.CinemaPhoto
	buckets[enum.TicketPDF] = env.Storage.Buckets.TicketPDF
	seenBuckets := make(map[string]enum.BucketType, len(buckets))
	for bucketType, bucketName := range buckets {
		bucketName = strings.TrimSpace(bucketName)
		if bucketName == "" {
			return nil, fmt.Errorf("bucket name for type %q is not configured", bucketType)
		}
		if existingType, duplicate := seenBuckets[bucketName]; duplicate {
			return nil, fmt.Errorf("bucket %q is configured for both %q and %q", bucketName, existingType, bucketType)
		}
		buckets[bucketType] = bucketName
		seenBuckets[bucketName] = bucketType
	}

	storage := &S3Storage{
		env:      env,
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		buckets:  buckets,
	}

	if err := storage.EnsureBuckets(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure storage buckets: %w", err)
	}

	return storage, nil
}

func (s *S3Storage) EnsureBuckets(ctx context.Context) error {
	for bucketType, bucketName := range s.buckets {
		if bucketName == "" {
			continue
		}

		_, err := s.client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			_, createErr := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if createErr != nil {
				return fmt.Errorf("failed to create bucket %q: %w", bucketName, createErr)
			}
		}

		// Ticket PDFs carry the admission QR credential. Keep their bucket
		// private and issue short-lived presigned URLs after authorization.
		if bucketType == enum.TicketPDF {
			continue
		}

		policy := fmt.Sprintf(`{
          "Version": "2012-10-17",
          "Statement": [
             {
                "Effect": "Allow",
                "Principal": {"AWS": ["*"]},
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
             }
          ]
       }`, bucketName)

		_, err = s.client.PutBucketPolicyWithContext(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucketName),
			Policy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("failed to set public policy for bucket %q: %w", bucketName, err)
		}
	}
	return nil
}

func (s *S3Storage) validateBucket(bucketType enum.BucketType) (string, error) {
	// تغییر مهم: استفاده مستقیم از متد IsValid به جای پکیج slices
	if !bucketType.IsValid() {
		return "", fmt.Errorf("bucket type %v does not exist", bucketType)
	}

	bucketName, ok := s.buckets[bucketType]
	if !ok || bucketName == "" {
		return "", fmt.Errorf("bucket name for type %v is not configured", bucketType)
	}

	return bucketName, nil
}

func (s *S3Storage) UploadObject(ctx context.Context, bucketType enum.BucketType, key string, file io.Reader, contentType string) (string, error) {
	bucketName, err := s.validateBucket(bucketType)
	if err != nil {
		return "", err
	}

	_, err = s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object %q to bucket %q: %w", key, bucketName, err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", s.env.Storage.Endpoint, bucketName, key)
	return fileURL, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, bucketType enum.BucketType, key string) error {
	bucketName, err := s.validateBucket(bucketType)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %q from bucket %q: %w", key, bucketName, err)
	}

	return nil
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, bucketType enum.BucketType, objectKey string, expiration time.Duration) (string, error) {
	bucketName, err := s.validateBucket(bucketType)
	if err != nil {
		return "", err
	}

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for %q: %w", objectKey, err)
	}

	return url, nil
}

func (s *S3Storage) GetObjectList(ctx context.Context, bucketType enum.BucketType) ([]string, error) {
	bucketName, err := s.validateBucket(bucketType)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %q: %w", bucketName, err)
	}

	var objects []string
	for _, item := range resp.Contents {
		if item.Key != nil {
			objects = append(objects, *item.Key)
		}
	}

	return objects, nil
}
