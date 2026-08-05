package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type RBACServiceInterface interface {
	// Assignments & Discharging
	AssignPermissionsToRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionID []*int64) error
	DischargePermissionsOfRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionID []*int64) error

	// Authorization Enforcement
	UserHasPermission(ctx context.Context, tokenData *dto.TokenInContext, userID int64, permissionName string) (bool, error)
	UserHasRole(ctx context.Context, tokenData *dto.TokenInContext, userID int64, roleName string) (bool, error)

	GetUserPermissions(ctx context.Context, tokenData *dto.TokenInContext, userID int64) ([]*dto.PermissionResponse, error)
	GetUserRoles(ctx context.Context, tokenData *dto.TokenInContext, userID int64) ([]*dto.RoleResponse, error)

	GetUsersByPermissionName(ctx context.Context, tokenData *dto.TokenInContext, permissionNames []string) ([]*dto.UserResponse, error)
	GetUsersByRoleNames(ctx context.Context, tokenData *dto.TokenInContext, roleNames []string) ([]*dto.UserResponse, error)

	GetPermissionsByRoleName(ctx context.Context, tokenData *dto.TokenInContext, role string) ([]*dto.PermissionResponse, error)
	GetRolesByPermissionName(ctx context.Context, tokenData *dto.TokenInContext, permissionName []string) ([]*dto.RoleResponse, error)

	GetPermissionsByID(ctx context.Context, tokenData *dto.TokenInContext, pid int64) (*dto.PermissionResponse, error)
	GetRoleByID(ctx context.Context, tokenData *dto.TokenInContext, rid int64) (*dto.RoleResponse, error)
	UpdateUserRoles(ctx context.Context, ctxToken *dto.TokenInContext, request dto.UpdateUserRolesRequest) (*dto.UserResponse, error)
}
