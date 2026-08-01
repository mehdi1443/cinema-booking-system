package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository"
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
			return nil, repository_contract.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return user, nil
}

func (r *userRepositoryPostgres) FindUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &user, nil
}

func (r *userRepositoryPostgres) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	res := r.db.GetDB().WithContext(ctx).Model(user).Where("id = ?", user.ID).Updates(user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, repository_contract.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, repository_contract.ErrUserNotFound
	}

	return user, nil
}
func (r *userRepositoryPostgres) CreateRole(ctx context.Context, role *entity.Role) error {
	err := r.db.GetDB().WithContext(ctx).Create(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository_contract.ErrRoleAlreadyExists
		}
		return fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrRoleNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &role, nil
}

func (r *userRepositoryPostgres) FindRoleByID(ctx context.Context, roleID int64) (*entity.Role, error) {
	var role entity.Role
	err := r.db.GetDB().WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrRoleNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &role, nil
}

func (r *userRepositoryPostgres) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	err := r.db.GetDB().WithContext(ctx).Create(permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository_contract.ErrPermissionAlreadyExists
		}
		return fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	var perm entity.Permission
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository_contract.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return &perm, nil
}

func (r *userRepositoryPostgres) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	var role entity.Role
	role.ID = roleID

	var perm entity.Permission
	perm.ID = permissionID

	err := r.db.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Append(&perm)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository_contract.ErrPermissionAlreadyExists
		}
		return fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error {
	var user entity.User
	user.ID = userID

	var role entity.Role
	role.ID = roleID

	err := r.db.GetDB().WithContext(ctx).Model(&user).Association("Roles").Append(&role)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository_contract.ErrRoleAlreadyExists
		}
		return fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *userRepositoryPostgres) UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *userRepositoryPostgres) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	var permissions []string
	err := r.db.GetDB().WithContext(ctx).
		Table("permissions").
		Select("DISTINCT permissions.name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository_contract.ErrDatabase, err)
	}
	return permissions, nil
}
