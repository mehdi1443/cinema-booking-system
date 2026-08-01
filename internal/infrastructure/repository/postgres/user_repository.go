package postgres

import (
	"context"
	"errors"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	repository "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/postgres"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

// import (
// 	"context"
// 	"errors"

// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/postgres"
// 	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
// 	"gorm.io/gorm"
// )

type userRepositoryPostgres struct{}

func NewUserRepositoryPostgres() repository.UserInterFace {
	return &userRepositoryPostgres{}
}

func (r *userRepositoryPostgres) CreateUser(ctx context.Context, db database.Database, user *entity.User) (*entity.User, error) {
	err := db.GetDB().WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryPostgres) FindUserByID(ctx context.Context, db database.Database, id int64) (*entity.User, error) {
	var user entity.User
	err := db.GetDB().WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByPhone(ctx context.Context, db database.Database, phone string) (*entity.User, error) {
	var user entity.User
	err := db.GetDB().WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPostgres) FindUserByEmail(ctx context.Context, db database.Database, email string) (*entity.User, error) {
	var user entity.User
	err := db.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPostgres) UpdateUser(ctx context.Context, db database.Database, user *entity.User) (*entity.User, error) {
	err := db.GetDB().WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryPostgres) CreateRole(ctx context.Context, db database.Database, role *entity.Role) error {
	return db.GetDB().WithContext(ctx).Create(role).Error
}

func (r *userRepositoryPostgres) FindRoleByName(ctx context.Context, db database.Database, name string) (*entity.Role, error) {
	var role entity.Role
	err := db.GetDB().WithContext(ctx).Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepositoryPostgres) FindRoleByID(ctx context.Context, db database.Database, roleID int64) (*entity.Role, error) {
	var role entity.Role
	err := db.GetDB().WithContext(ctx).First(&role, roleID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepositoryPostgres) CreatePermission(ctx context.Context, db database.Database, permission *entity.Permission) error {
	return db.GetDB().WithContext(ctx).Create(permission).Error
}

func (r *userRepositoryPostgres) FindPermissionByName(ctx context.Context, db database.Database, name string) (*entity.Permission, error) {
	var perm entity.Permission
	err := db.GetDB().WithContext(ctx).Where("name = ?", name).First(&perm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *userRepositoryPostgres) AssignPermissionToRole(ctx context.Context, db database.Database, roleID int64, permissionID int64) error {
	var role entity.Role
	role.ID = roleID

	var perm entity.Permission
	perm.ID = permissionID

	return db.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Append(&perm)
}

func (r *userRepositoryPostgres) AssignRoleToUser(ctx context.Context, db database.Database, userID int64, roleID int64) error {
	var user entity.User
	user.ID = userID

	var role entity.Role
	role.ID = roleID

	return db.GetDB().WithContext(ctx).Model(&user).Association("Roles").Append(&role)
}

func (r *userRepositoryPostgres) RoleHasPermission(ctx context.Context, db database.Database, roleID int64, permissionID int64) (bool, error) {
	var count int64
	err := db.GetDB().WithContext(ctx).
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryPostgres) UserHasRole(ctx context.Context, db database.Database, userID int64, roleID int64) (bool, error) {
	var count int64
	err := db.GetDB().WithContext(ctx).
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryPostgres) GetUserPermissions(ctx context.Context, db database.Database, userID int64) ([]string, error) {
	var permissions []string
	err := db.GetDB().WithContext(ctx).
		Table("permissions").
		Select("DISTINCT permissions.name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}
	return permissions, nil
}
