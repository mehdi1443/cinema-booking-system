package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"

	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type userRepositoryPostgres struct {
	db database.Database
}

func NewUserRepositoryPostgres(db database.Database) repository_contract.UserInterFace {
	return &userRepositoryPostgres{
		db: db,
	}
}

func (r *userRepositoryPostgres) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.db.GetDB().WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return user, nil
}

func (r *userRepositoryPostgres) FindUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	res := r.db.GetDB().WithContext(ctx).Model(user).Where("id = ?", user.ID).Updates(user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrUserNotFound
	}

	return user, nil
}

func (r *userRepositoryPostgres) FindUsersByStatus(ctx context.Context, status enum.UserStatus) ([]*entity.User, error) {
	var users []*entity.User

	err := r.db.GetDB().WithContext(ctx).
		Preload("Roles").
		Where("status = ?", status).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return users, nil
}

func (r *userRepositoryPostgres) DeleteUserByID(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryPostgres) FindUsers(ctx context.Context, filter *dto.GetUsersQueryDTO) ([]*entity.User, error) {
	var users []*entity.User
	query := r.db.GetDB().WithContext(ctx).Model(&entity.User{})

	if filter.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filter.Email+"%")
	}
	if filter.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+filter.LastName+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if len(filter.Roles) > 0 {
		query = query.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.name IN ?", filter.Roles)
	}
	if len(filter.Permissions) > 0 {
		query = query.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("permissions.name IN ?", filter.Permissions)
	}

	err := query.Distinct().Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return users, nil
}

func (r *userRepositoryPostgres) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	res := r.db.GetDB().WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Update("password", hashedPassword)

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryPostgres) VerifyUserPhone(ctx context.Context, userID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"is_phone_verified": true,
			"phone_verified_at": time.Now(),
		})

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}
	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}
	return nil
}

func (r *userRepositoryPostgres) VerifyUserEmail(ctx context.Context, userID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"is_email_verified": true,
			"email_verified_at": time.Now(),
		})

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}
	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}
	return nil
}
