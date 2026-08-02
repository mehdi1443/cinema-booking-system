package redis_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type UserPermissionCacheRepository struct {
	client database.Cache
}

func NewUserPermissionCacheRepository(client database.Cache) *UserPermissionCacheRepository {
	return &UserPermissionCacheRepository{client: client}
}

func (r *UserPermissionCacheRepository) getkeysignuture(userID string) string {
	// user:user_id:permissions : REDIS_SET_DATATYPE
	return fmt.Sprintf("%s:%s:%s", r.client.GetRedisConstant().UserSignature, userID, r.client.GetRedisConstant().UserPermissionSignature)
}

func (r *UserPermissionCacheRepository) SaveUserPermissions(ctx context.Context, userID string, permissions []string, ttl time.Duration) error {
	key := r.getkeysignuture(userID)

	if err := r.client.GetRDB().SAdd(ctx, key, permissions, ttl).Err(); err != nil {
		return fmt.Errorf("redis save user permissions error: %w", err)
	}

	return nil
}

func (r *UserPermissionCacheRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := r.getkeysignuture(userID)

	permissions, err := r.client.GetRDB().SMembers(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repository_contract.ErrUserPermissionNotFound
		}
		return nil, fmt.Errorf("redis get user permissions error: %w", err)
	}

	return permissions, nil
}

func (r *UserPermissionCacheRepository) DeleteUserPermission(ctx context.Context, userID string) error {
	key := r.getkeysignuture(userID)

	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user permission error: %w", err)
	}

	return nil
}

func (r *UserPermissionCacheRepository) DeleteAllUserPermission(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = r.getkeysignuture(id)
	}

	if err := r.client.GetRDB().Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete all user permissions error: %w", err)
	}

	return nil
}

func (r *UserPermissionCacheRepository) IsExistsUserPermissions(ctx context.Context, userID string) (*bool, error) {
	key := r.getkeysignuture(userID)

	count, err := r.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis exists user permissions error: %w", err)
	}

	exists := count > 0
	return &exists, nil
}
