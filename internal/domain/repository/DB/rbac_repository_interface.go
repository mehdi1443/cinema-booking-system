package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type RBACInterface interface {
	FindPermissionByID(ctx context.Context, pid int64) (*entity.Permission, error)
	CreateRole(ctx context.Context, role *entity.Role) error
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
	FindRoleByID(ctx context.Context, roleID int64) (*entity.Role, error)
	CreatePermission(ctx context.Context, permission *entity.Permission) error // just used for seed!
	FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error)
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

	UpdateUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
	AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error
	UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error)
}
