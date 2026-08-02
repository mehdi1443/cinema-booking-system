package cache

import (
	"context"
	"errors"
	"time"
)

var ErrUserPermissionNotFound = errors.New("user permission not found in redis")

type UserPermissionCacheInterface interface {
	SaveUserPermissions(ctx context.Context, userID string, permission []string, ttl time.Duration) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	DeleteUserPermission(ctx context.Context, userID string) error
	DeleteAllUserPermission(ctx context.Context, userID []string) error
	IsExistsUserPermissions(ctx context.Context, userID string) (*bool, error)
}
