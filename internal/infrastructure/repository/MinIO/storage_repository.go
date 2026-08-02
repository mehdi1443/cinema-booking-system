package MinIO

import (
	"context"
	"fmt"
	"io"
	"slices"
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

	return &S3Storage{
		env:      env,
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		buckets:  buckets,
	}, nil
}

func (s *S3Storage) validateBucket(bucketType enum.BucketType) (string, error) {
	validTypes := enum.GetAllBucketTypes()
	if !slices.Contains(validTypes, bucketType) {
		return "", fmt.Errorf("bucket type %v does not exist", bucketType)
	}

	bucketName, ok := s.buckets[bucketType]
	if !ok || bucketName == "" {
		return "", fmt.Errorf("bucket name for type %v is not configured", bucketType)
	}

	return bucketName, nil
}

func (s *S3Storage) UploadObject(ctx context.Context, bucketType enum.BucketType, key string, file io.Reader) (string, error) {
	bucketName, err := s.validateBucket(bucketType)
	if err != nil {
		return "", err
	}

	_, err = s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
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
