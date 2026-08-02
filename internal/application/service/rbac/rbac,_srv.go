package rbacsrv

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type rbacService struct {
	repo repository_contract.UserInterFace
}

func NewRBACService(repo repository_contract.UserInterFace) service_contract.RBACServiceInterface {
	return &rbacService{repo: repo}
}

func (s *rbacService) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return s.repo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) DischargePermissionsOfRole(ctx context.Context, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return s.repo.DischargePermissionsOfRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) UserHasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	return s.repo.UserHasPermission(ctx, userID, permissionName)
}

func (s *rbacService) UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	return s.repo.UserHasRoleByName(ctx, userID, roleName)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID int64) ([]*dto.PermissionResponse, error) {
	perms, err := s.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseSlice(perms), nil
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID int64) ([]*dto.RoleResponse, error) {
	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponseSlice(roles), nil
}

func (s *rbacService) GetUsersByPermissionName(ctx context.Context, permissionNames []string) ([]*dto.UserResponse, error) {
	users, err := s.repo.GetUsersByPermissionNames(ctx, permissionNames)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponseSlice(users, &dto.HideFieldUserOption{WalletBalance: false}), nil
}

func (s *rbacService) GetUsersByRoleNames(ctx context.Context, roleNames []string) ([]*dto.UserResponse, error) {
	users, err := s.repo.GetUsersByRoleNames(ctx, roleNames)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponseSlice(users, &dto.HideFieldUserOption{WalletBalance: false}), nil
}

func (s *rbacService) GetPermissionsByRoleName(ctx context.Context, role string) ([]*dto.PermissionResponse, error) {
	perms, err := s.repo.GetPermissionsByRoleName(ctx, role)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseSlice(perms), nil
}

func (s *rbacService) GetRolesByPermissionName(ctx context.Context, permissionName []string) ([]*dto.RoleResponse, error) {
	roles, err := s.repo.GetRolesByPermissionNames(ctx, permissionName)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponseSlice(roles), nil
}

func (s *rbacService) GetPermissionsByID(ctx context.Context, pid int64) (*dto.PermissionResponse, error) {
	perm, err := s.repo.FindPermissionByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponse(perm), nil
}

func (s *rbacService) GetRoleByID(ctx context.Context, rid int64) (*dto.RoleResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, rid)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponse(role), nil
}
