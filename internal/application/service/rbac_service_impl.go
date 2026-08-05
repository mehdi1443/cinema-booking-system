package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type rbacService struct {
	rbacRepo repository_contract.RBACInterface
	userRepo repository_contract.UserInterFace
}

func NewRBACService(rbacRepo repository_contract.RBACInterface, userRepo repository_contract.UserInterFace) service_contract.RBACServiceInterface {
	return &rbacService{rbacRepo: rbacRepo, userRepo: userRepo}
}

func (s *rbacService) AssignPermissionsToRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return s.rbacRepo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) DischargePermissionsOfRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return s.rbacRepo.DischargePermissionsOfRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) UserHasPermission(ctx context.Context, tokenData *dto.TokenInContext, userID int64, permissionName string) (bool, error) {
	return s.rbacRepo.UserHasPermission(ctx, userID, permissionName)
}

func (s *rbacService) UserHasRole(ctx context.Context, tokenData *dto.TokenInContext, userID int64, roleName string) (bool, error) {
	return s.rbacRepo.UserHasRoleByName(ctx, userID, roleName)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, tokenData *dto.TokenInContext, userID int64) ([]*dto.PermissionResponse, error) {
	perms, err := s.rbacRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseSlice(perms), nil
}

func (s *rbacService) GetUserRoles(ctx context.Context, tokenData *dto.TokenInContext, userID int64) ([]*dto.RoleResponse, error) {
	roles, err := s.rbacRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponseSlice(roles), nil
}

func (s *rbacService) GetUsersByPermissionName(ctx context.Context, tokenData *dto.TokenInContext, permissionNames []string) ([]*dto.UserResponse, error) {
	users, err := s.rbacRepo.GetUsersByPermissionNames(ctx, permissionNames)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponseSlice(users, &dto.HideFieldUserOption{WalletBalance: false}), nil
}

func (s *rbacService) GetUsersByRoleNames(ctx context.Context, tokenData *dto.TokenInContext, roleNames []string) ([]*dto.UserResponse, error) {
	users, err := s.rbacRepo.GetUsersByRoleNames(ctx, roleNames)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponseSlice(users, &dto.HideFieldUserOption{WalletBalance: false}), nil
}

func (s *rbacService) GetPermissionsByRoleName(ctx context.Context, tokenData *dto.TokenInContext, role string) ([]*dto.PermissionResponse, error) {
	perms, err := s.rbacRepo.GetPermissionsByRoleName(ctx, role)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponseSlice(perms), nil
}

func (s *rbacService) GetRolesByPermissionName(ctx context.Context, tokenData *dto.TokenInContext, permissionName []string) ([]*dto.RoleResponse, error) {
	roles, err := s.rbacRepo.GetRolesByPermissionNames(ctx, permissionName)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponseSlice(roles), nil
}

func (s *rbacService) GetPermissionsByID(ctx context.Context, tokenData *dto.TokenInContext, pid int64) (*dto.PermissionResponse, error) {
	perm, err := s.rbacRepo.FindPermissionByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	return dto.ToPermissionResponse(perm), nil
}

func (s *rbacService) GetRoleByID(ctx context.Context, tokenData *dto.TokenInContext, rid int64) (*dto.RoleResponse, error) {
	role, err := s.rbacRepo.FindRoleByID(ctx, rid)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponse(role), nil
}

func (s *rbacService) UpdateUserRoles(ctx context.Context, ctxToken *dto.TokenInContext, req dto.UpdateUserRolesRequest) (*dto.UserResponse, error) {
	// 1. Verify user exists
	_, err := s.userRepo.FindUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 2. Sync roles via repository
	err = s.rbacRepo.UpdateUserRoles(ctx, req.UserID, req.RoleIDs)
	if err != nil {
		return nil, err
	}

	// 3. Fetch reloaded user with new roles
	updatedUser, err := s.userRepo.FindUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(updatedUser, &dto.HideFieldUserOption{WalletBalance: false}), nil
}
