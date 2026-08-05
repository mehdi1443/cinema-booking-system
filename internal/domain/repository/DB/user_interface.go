package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type UserInterFace interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByID(ctx context.Context, id int64) (*entity.User, error)
	FindUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)

	CreateRole(ctx context.Context, role *entity.Role) error
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
	FindRoleByID(ctx context.Context, roleID int64) (*entity.Role, error)

	CreatePermission(ctx context.Context, permission *entity.Permission) error // just used for seed!
	FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error)

	FindUsersByStatus(ctx context.Context, status enum.UserStatus) ([]*entity.User, error)
	DeleteUserByID(ctx context.Context, id int64) error
	UpdateUserRoles(ctx context.Context, userID int64, roleIDs []int64) error

	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
	AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error
	RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error)
	UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error)

	AddAmountToBalance(ctx context.Context, userID int64, amount int64) (int64, error)

	AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []*int64) error

	DischargePermissionsOfRole(ctx context.Context, roleID int64, permissionIDs []*int64) error

	UserHasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)

	UserHasRoleByName(ctx context.Context, userID int64, roleName string) (bool, error)

	GetUserPermissions(ctx context.Context, userID int64) ([]*entity.Permission, error)

	GetUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error)

	GetUsersByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.User, error)

	GetUsersByRoleNames(ctx context.Context, roleNames []string) ([]*entity.User, error)

	GetPermissionsByRoleName(ctx context.Context, roleName string) ([]*entity.Permission, error)

	GetRolesByPermissionNames(ctx context.Context, permissionNames []string) ([]*entity.Role, error)

	FindPermissionByID(ctx context.Context, pid int64) (*entity.Permission, error)

	FindUsers(ctx context.Context, filter *dto.GetUsersQueryDTO) ([]*entity.User, error)

	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error

	VerifyUserPhone(ctx context.Context, userID int64) error
	VerifyUserEmail(ctx context.Context, userID int64) error
}
