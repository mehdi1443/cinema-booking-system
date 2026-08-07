package storage

import (
	"context"
	"io"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type StorageInterFace interface {
	UploadObject(ctx context.Context, bucketType enum.BucketType, key string, file io.Reader, contentType string) (string, error)
	DeleteObject(ctx context.Context, bucketType enum.BucketType, key string) error
	GetPresignedURL(ctx context.Context, bucketType enum.BucketType, objectKey string, expiration time.Duration) (string, error)
	GetObjectList(ctx context.Context, bucketType enum.BucketType) ([]string, error)
}
