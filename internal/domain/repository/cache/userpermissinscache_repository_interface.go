package cache

import (
	"context"
	"time"
)

type UserPermissionCacheInterface interface {
	SaveUserPermissions(ctx context.Context, userID string, permission []string, ttl time.Duration) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	DeleteUserPermission(ctx context.Context, userID string) error
	DeleteAllUserPermission(ctx context.Context, userID []string) error
	IsExistsUserPermissions(ctx context.Context, userID string) (*bool, error)
}
