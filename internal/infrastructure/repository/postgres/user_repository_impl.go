package postgres_repo

import (
	"context"
	"errors"
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

func NewUserRepositoryPostgres(db database.Database) repository_contract.UserInterface {
	return &userRepositoryPostgres{
		db: db,
	}
}

func (r *userRepositoryPostgres) CreateUser(
	ctx context.Context,
	user *entity.User,
) (*entity.User, error) {

	err := r.db.
		GetDB().
		WithContext(ctx).
		Create(user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrUserAlreadyExists
		}
		return nil, handleDBError(err)
	}

	return user, nil
}

func (r *userRepositoryPostgres) CreateUserWithRole(
	ctx context.Context,
	user *entity.User,
	roleName string,
) (*entity.User, error) {
	err := r.db.WithTransaction(func(txDatabase database.Database) error {
		tx := txDatabase.GetDB().WithContext(ctx)
		var role entity.Role
		if err := tx.Where("name = ? AND deleted_at IS NULL", roleName).First(&role).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrRoleNotFound
			}
			return handleDBError(err)
		}
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return exception.ErrUserAlreadyExists
			}
			return handleDBError(err)
		}
		if err := tx.Model(user).Association("Roles").Append(&role); err != nil {
			return handleDBError(err)
		}
		user.Roles = []*entity.Role{&role}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryPostgres) FindUserByID(
	ctx context.Context,
	id int64,
) (*entity.User, error) {

	var user entity.User

	err := r.db.
		GetDB().
		WithContext(ctx).
		Preload("Roles").
		First(&user, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, handleDBError(err)
	}

	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByPhone(
	ctx context.Context,
	phone string,
) (*entity.User, error) {

	var user entity.User

	phones := []string{phone}
	if len(phone) == 11 && phone[:2] == "09" {
		phones = append(phones, phone[1:], "+98"+phone[1:])
	}

	err := r.db.
		GetDB().
		WithContext(ctx).
		Preload("Roles").
		Where("phone IN ?", phones).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, handleDBError(err)
	}

	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {

	var user entity.User

	err := r.db.
		GetDB().
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, handleDBError(err)
	}

	return &user, nil
}

func (r *userRepositoryPostgres) UpdateUser(
	ctx context.Context,
	user *entity.User,
) (*entity.User, error) {

	res := r.db.
		GetDB().
		WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", user.ID).
		// This method is used for profile/status updates. Updating the whole
		// loaded entity would overwrite a concurrently changed wallet balance,
		// password, roles, or verification state with stale values.
		Updates(map[string]any{
			"first_name":        user.FirstName,
			"last_name":         user.LastName,
			"email":             user.Email,
			"is_email_verified": user.IsEmailVerified,
			"email_verified_at": user.EmailVerifiedAt,
			"status":            user.Status,
			"updated_at":        time.Now(),
		})

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrUserAlreadyExists
		}
		return nil, handleDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrUserNotFound
	}

	return user, nil
}

func (r *userRepositoryPostgres) FindUsersByStatus(
	ctx context.Context,
	status enum.UserStatus,
) ([]*entity.User, error) {

	var users []*entity.User

	err := r.db.
		GetDB().
		WithContext(ctx).
		Preload("Roles").
		Where("status = ?", status).
		Find(&users).
		Error

	if err != nil {
		return nil, handleDBError(err)
	}

	return users, nil
}

func (r *userRepositoryPostgres) DeleteUserByID(
	ctx context.Context,
	id int64,
) error {

	res := r.db.
		GetDB().
		WithContext(ctx).
		Delete(&entity.User{}, id)

	if res.Error != nil {
		return handleDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryPostgres) FindUsers(
	ctx context.Context,
	filter *dto.GetUsersQueryDTO,
) ([]*entity.User, error) {

	var users []*entity.User

	query := r.db.
		GetDB().
		WithContext(ctx).
		Model(&entity.User{})

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
		query = query.
			Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.name IN ?", filter.Roles)
	}

	if len(filter.Permissions) > 0 {
		query = query.
			Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("permissions.name IN ?", filter.Permissions)
	}

	err := query.
		Distinct().
		Find(&users).
		Error

	if err != nil {
		return nil, handleDBError(err)
	}

	return users, nil
}

func (r *userRepositoryPostgres) UpdatePassword(
	ctx context.Context,
	userID int64,
	hashedPassword string,
) error {

	res := r.db.
		GetDB().
		WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Update("password_hash", hashedPassword)

	if res.Error != nil {
		return handleDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryPostgres) VerifyUserPhone(
	ctx context.Context,
	userID int64,
) error {

	res := r.db.
		GetDB().
		WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"is_phone_verified": true,
			"phone_verified_at": time.Now(),
		})

	if res.Error != nil {
		return handleDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryPostgres) VerifyUserEmail(
	ctx context.Context,
	email string,
) error {

	res := r.db.
		GetDB().
		WithContext(ctx).
		Model(&entity.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Updates(map[string]interface{}{
			"is_email_verified": true,
			"email_verified_at": time.Now(),
		})

	if res.Error != nil {
		return handleDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}
