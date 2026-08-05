package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
)

type rbacService struct {
	rbacRepo                  repository_contract.RBACInterface
	userRepo                  repository_contract.UserInterface
	userPermissionCacheRepo   cache.UserPermissionCacheInterface
	userPermissionCacheConfig *bootstrap.UserPermissionCache
}

func NewRBACService(
	rbacRepo repository_contract.RBACInterface,
	userRepo repository_contract.UserInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	userPermissionCacheConfig *bootstrap.UserPermissionCache,
) service_contract.RBACServiceInterface {
	return &rbacService{rbacRepo: rbacRepo, userRepo: userRepo,
		userPermissionCacheRepo: userPermissionCacheRepo, userPermissionCacheConfig: userPermissionCacheConfig}
}

func (s *rbacService) AssignPermissionsToRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	err := s.rbacRepo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
	if err != nil {
		return err
	}

	role, err := s.rbacRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	users, err := s.GetUsersByRoleNames(ctx, tokenData, []string{role.Name})
	if err != nil {
		return err
	}
	userIDs := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		id := strconv.Itoa(int(users[i].ID))
		userIDs = append(userIDs, id)
	}
	err = s.userPermissionCacheRepo.DeleteAllUserPermission(ctx, userIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *rbacService) DischargePermissionsOfRole(ctx context.Context, tokenData *dto.TokenInContext, roleID int64, permissionIDs []*int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	err := s.rbacRepo.DischargePermissionsOfRole(ctx, roleID, permissionIDs)
	if err != nil {
		return err
	}

	role, err := s.rbacRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	users, err := s.GetUsersByRoleNames(ctx, tokenData, []string{role.Name})
	if err != nil {
		return err
	}
	userIDs := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		id := strconv.Itoa(int(users[i].ID))
		userIDs = append(userIDs, id)
	}
	err = s.userPermissionCacheRepo.DeleteAllUserPermission(ctx, userIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *rbacService) UserHasPermission(ctx context.Context, tokenData *dto.TokenInContext, userID int64, permissionName string) (bool, error) {
	return s.rbacRepo.UserHasPermission(ctx, userID, permissionName)
}

func (s *rbacService) UserHasRole(ctx context.Context, tokenData *dto.TokenInContext, userID int64, roleName string) (bool, error) {
	return s.rbacRepo.UserHasRoleByName(ctx, userID, roleName)
}

func (s *rbacService) GetUserPermissionNames(ctx context.Context, tokenData *dto.TokenInContext, userID int64) ([]string, error) {
	permsStr, err := s.userPermissionCacheRepo.GetUserPermissions(ctx, strconv.Itoa(int(userID)))
	if err != nil {
		if errors.Is(err, exception.ErrUserPermissionNotFound) {
			permsStr, err = s.rbacRepo.GetUserPermissionNames(ctx, userID)
			if err != nil {
				return nil, err
			}
			err := s.userPermissionCacheRepo.SaveUserPermissions(ctx, strconv.Itoa(int(userID)), permsStr, time.Duration(s.userPermissionCacheConfig.UserPermissionExpiredDuration)*time.Minute)
			if err != nil {
				return nil, err
			}
			return permsStr, nil
		}
		return nil, err

	}
	return permsStr, nil

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

func (s *rbacService) UpdateUserRoles(ctx context.Context, ctxToken *dto.TokenInContext, userID int64, req dto.UpdateUserRolesRequest) (*dto.UserResponse, error) {
	// 1. Verify user exists
	_, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Sync roles via repository
	err = s.rbacRepo.UpdateUserRoles(ctx, userID, req.RoleIDs)
	if err != nil {
		return nil, err
	}
	err = s.userPermissionCacheRepo.DeleteUserPermission(ctx, strconv.Itoa(int(userID)))
	if err != nil {
		return nil, err
	}

	// 3. Fetch reloaded user with new roles
	updatedUser, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(updatedUser, &dto.HideFieldUserOption{WalletBalance: false}), nil
}
