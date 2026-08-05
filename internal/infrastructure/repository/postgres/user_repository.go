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
	"gorm.io/gorm/clause"
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
func (r *userRepositoryPostgres) CreateRole(ctx context.Context, role *entity.Role) error {
	err := r.db.GetDB().WithContext(ctx).Create(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrRoleAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
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

func (r *userRepositoryPostgres) FindRoleByID(ctx context.Context, roleID int64) (*entity.Role, error) {
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

func (r *userRepositoryPostgres) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	err := r.db.GetDB().WithContext(ctx).Create(permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrPermissionAlreadyExists
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *userRepositoryPostgres) FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
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

func (r *userRepositoryPostgres) AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error {
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

func (r *userRepositoryPostgres) AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error {
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

func (r *userRepositoryPostgres) RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error) {
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

func (r *userRepositoryPostgres) UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error) {
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

// DecreaseBalance atomically deducts funds inside a database transaction with explicit row locking.
func (r *userRepositoryPostgres) AddAmountToBalance(ctx context.Context, userID int64, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, exception.ErrInvalidAmount
	}

	var finalBalance int64

	// Begin database transaction block
	err := r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user entity.User

		//  Lock the user row (FOR UPDATE) to prevent race conditions from concurrent requests
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", userID).
			First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrUserNotFound
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if user.WalletBalance+amount < 0 {
			return exception.ErrInsufficientBalance
		} else if user.WalletBalance+amount > entity.MaxWalletBalance {
			return exception.ErrReachManBalance
		}

		finalBalance = user.WalletBalance + amount
		err = tx.Model(&user).Update("wallet_balance", finalBalance).Error
		if err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return finalBalance, nil
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

func (r *userRepositoryPostgres) UpdateUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
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

func (r *userRepositoryPostgres) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
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

func (r *userRepositoryPostgres) DischargePermissionsOfRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
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

func (r *userRepositoryPostgres) UserHasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
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

func (r *userRepositoryPostgres) UserHasRoleByName(ctx context.Context, userID int64, roleName string) (bool, error) {
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

func (r *userRepositoryPostgres) GetUserPermissions(ctx context.Context, userID int64) ([]*entity.Permission, error) {
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

func (r *userRepositoryPostgres) GetUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error) {
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

func (r *userRepositoryPostgres) GetUsersByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.User, error) {
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

func (r *userRepositoryPostgres) GetUsersByRoleNames(ctx context.Context, roleNames []string) ([]*entity.User, error) {
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

func (r *userRepositoryPostgres) GetPermissionsByRoleName(ctx context.Context, roleName string) ([]*entity.Permission, error) {
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

func (r *userRepositoryPostgres) GetRolesByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.Role, error) {
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

func (r *userRepositoryPostgres) FindPermissionByID(ctx context.Context, pid int64) (*entity.Permission, error) {
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
