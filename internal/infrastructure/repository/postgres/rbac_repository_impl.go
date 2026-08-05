package postgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type rbacRepositoryPostgres struct {
	db database.Database
}

func NewRABCRepositoryPostgres(db database.Database) repository_contract.RBACInterface {
	return &rbacRepositoryPostgres{
		db: db,
	}
}

func (r *rbacRepositoryPostgres) GetRolesByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.GetDB().WithContext(ctx).
		Distinct("roles.*").
		Table("roles").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("permissions.name IN ? AND roles.deleted_at IS NULL", permissionNames).
		Preload("Permissions").
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return roles, nil
}

func (r *rbacRepositoryPostgres) FindPermissionByID(ctx context.Context, pid int64) (*entity.Permission, error) {
	var perm entity.Permission
	err := r.db.GetDB().WithContext(ctx).First(&perm, pid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &perm, nil
}

func (r *rbacRepositoryPostgres) CreateRole(ctx context.Context, role *entity.Role) error {
	err := r.db.GetDB().WithContext(ctx).Create(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrRoleAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrRoleNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &role, nil
}

func (r *rbacRepositoryPostgres) FindRoleByID(ctx context.Context, roleID int64) (*entity.Role, error) {
	var role entity.Role
	err := r.db.GetDB().WithContext(ctx).First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrRoleNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &role, nil
}

func (r *rbacRepositoryPostgres) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	err := r.db.GetDB().WithContext(ctx).Create(permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrPermissionAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	var perm entity.Permission
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &perm, nil
}

func (r *rbacRepositoryPostgres) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
	var role entity.Role
	role.ID = roleID

	var perm entity.Permission
	perm.ID = permissionID

	err := r.db.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Append(&perm)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrPermissionAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error {
	var user entity.User
	user.ID = userID

	var role entity.Role
	role.ID = roleID

	err := r.db.GetDB().WithContext(ctx).Model(&user).Association("Roles").Append(&role)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrRoleAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *rbacRepositoryPostgres) UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *rbacRepositoryPostgres) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
	role := entity.Role{BaseEntity: entity.BaseEntity{ID: roleID}}
	perms := make([]*entity.Permission, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		perms = append(perms, &entity.Permission{BaseEntity: entity.BaseEntity{ID: *pid}})
	}

	err := r.db.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Append(perms)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) DischargePermissionsOfRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
	role := entity.Role{BaseEntity: entity.BaseEntity{ID: roleID}}
	perms := make([]*entity.Permission, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		perms = append(perms, &entity.Permission{BaseEntity: entity.BaseEntity{ID: *pid}})
	}

	err := r.db.GetDB().WithContext(ctx).Model(&role).Association("Permissions").Delete(perms)
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *rbacRepositoryPostgres) UserHasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.name = ? AND permissions.deleted_at IS NULL", userID, permissionName).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *rbacRepositoryPostgres) UserHasRoleByName(ctx context.Context, userID int64, roleName string) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ? AND roles.deleted_at IS NULL", userID, roleName).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return count > 0, nil
}

func (r *rbacRepositoryPostgres) GetUserPermissions(ctx context.Context, userID int64) ([]*entity.Permission, error) {
	var perms []*entity.Permission
	err := r.db.GetDB().WithContext(ctx).
		Distinct("permissions.*").
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL", userID).
		Find(&perms).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return perms, nil
}

func (r *rbacRepositoryPostgres) GetUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.GetDB().WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", userID).
		Preload("Permissions").
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return roles, nil
}

func (r *rbacRepositoryPostgres) GetUsersByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.GetDB().WithContext(ctx).
		Distinct("users.*").
		Table("users").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("permissions.name IN ? AND users.deleted_at IS NULL", permissionNames).
		Preload("Roles").
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return users, nil
}

func (r *rbacRepositoryPostgres) GetUsersByRoleNames(ctx context.Context, roleNames []string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.GetDB().WithContext(ctx).
		Distinct("users.*").
		Table("users").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name IN ? AND users.deleted_at IS NULL", roleNames).
		Preload("Roles").
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return users, nil
}

func (r *rbacRepositoryPostgres) GetPermissionsByRoleName(ctx context.Context, roleName string) ([]*entity.Permission, error) {
	var perms []*entity.Permission
	err := r.db.GetDB().WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Where("roles.name = ? AND permissions.deleted_at IS NULL", roleName).
		Find(&perms).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return perms, nil
}

func (r *rbacRepositoryPostgres) UpdateUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	user := entity.User{
		BaseEntity: entity.BaseEntity{
			ID: userID,
		},
	}

	roles := make([]*entity.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		roles = append(roles, &entity.Role{BaseEntity: entity.BaseEntity{
			ID: id,
		}})
	}

	err := r.db.GetDB().WithContext(ctx).
		Model(&user).
		Association("Roles").
		Replace(roles)

	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return nil
}
