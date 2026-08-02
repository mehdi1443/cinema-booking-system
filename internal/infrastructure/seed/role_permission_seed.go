package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
)

func strPtr(s string) *string {
	return &s
}

type RoleSeeder struct {
	superAdmin bootstrap.SuperUser
	userRepo   repository_contract.UserInterFace
	db         database.Database
	hasher     *hash.BcryptHasher
}

func NewRoleSeeder(
	superAdmin bootstrap.SuperUser,
	userRepo repository_contract.UserInterFace,
	db database.Database,
	hasher *hash.BcryptHasher,
) *RoleSeeder {
	return &RoleSeeder{
		superAdmin: superAdmin,
		userRepo:   userRepo,
		db:         db,
		hasher:     hasher,
	}
}

func (s *RoleSeeder) Seed(ctx context.Context) error {
	permMap, err := s.seedPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	roleMap, err := s.seedRolesWithPermissions(ctx, permMap)
	if err != nil {
		return fmt.Errorf("failed to seed roles and permissions: %w", err)
	}

	if err := s.seedSuperAdmin(ctx, roleMap[enum.RoleAdmin]); err != nil {
		return fmt.Errorf("failed to seed super admin: %w", err)
	}

	return nil
}

func (s *RoleSeeder) seedPermissions(ctx context.Context) (map[enum.PermissionType]*entity.Permission, error) {
	permMap := make(map[enum.PermissionType]*entity.Permission)

	for _, permType := range enum.GetAllPermissionTypes() {
		permName := permType.String()
		perm, err := s.userRepo.FindPermissionByName(ctx, permName)
		if err != nil && !errors.Is(err, repository_contract.ErrPermissionNotFound) {
			return nil, err
		}

		if errors.Is(err, repository_contract.ErrPermissionNotFound) {
			perm = &entity.Permission{
				Name:        permName,
				Title:       permType.Description(),
				Category:    permType.Category(),
				Description: permType.Description(),
			}
			if err := s.userRepo.CreatePermission(ctx, perm); err != nil {
				return nil, err
			}
		}
		permMap[permType] = perm
	}

	return permMap, nil
}

func (s *RoleSeeder) seedRolesWithPermissions(
	ctx context.Context,
	permMap map[enum.PermissionType]*entity.Permission,
) (map[enum.RoleName]*entity.Role, error) {
	roleMap := make(map[enum.RoleName]*entity.Role)

	for _, roleEnum := range enum.GetAllRoles() {
		role, err := s.userRepo.FindRoleByName(ctx, roleEnum.String())
		if err != nil && !errors.Is(err, repository_contract.ErrRoleNotFound) {
			return nil, err
		}

		if errors.Is(err, repository_contract.ErrRoleNotFound) {
			role = &entity.Role{
				Name:  roleEnum.String(),
				Title: roleEnum.Title(),
			}
			if err := s.userRepo.CreateRole(ctx, role); err != nil {
				return nil, err
			}
		}
		roleMap[roleEnum] = role

		for _, permType := range roleEnum.DefaultPermissions() {
			if permType == enum.PermissionAll {
				for _, perm := range permMap {
					if err := s.assignPermissionIfNotExist(ctx, role.ID, perm.ID); err != nil {
						return nil, err
					}
				}
			} else if perm, exists := permMap[permType]; exists {
				if err := s.assignPermissionIfNotExist(ctx, role.ID, perm.ID); err != nil {
					return nil, err
				}
			}
		}
	}

	return roleMap, nil
}

func (s *RoleSeeder) assignPermissionIfNotExist(ctx context.Context, roleID, permID int64) error {
	hasPerm, err := s.userRepo.RoleHasPermission(ctx, roleID, permID)
	if err != nil {
		return err
	}
	if !hasPerm {
		return s.userRepo.AssignPermissionToRole(ctx, roleID, permID)
	}
	return nil
}

func (s *RoleSeeder) seedSuperAdmin(ctx context.Context, adminRole *entity.Role) error {
	user, err := s.userRepo.FindUserByPhone(ctx, s.superAdmin.Phone)
	if err != nil && !errors.Is(err, repository_contract.ErrUserNotFound) {
		return err
	}

	if errors.Is(err, repository_contract.ErrUserNotFound) {
		hashedPassword, err := s.hasher.Hash(s.superAdmin.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = &entity.User{
			Phone:         s.superAdmin.Phone,
			Email:         s.superAdmin.Email,
			PasswordHash:  hashedPassword,
			FirstName:     s.superAdmin.FirstName,
			LastName:      s.superAdmin.LastName,
			WalletBalance: s.superAdmin.Wallet,
			Status:        "ACTIVE",
		}

		user, err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			return err
		}
	}

	hasRole, err := s.userRepo.UserHasRole(ctx, user.ID, adminRole.ID)
	if err != nil {
		return err
	}
	if !hasRole {
		if err := s.userRepo.AssignRoleToUser(ctx, user.ID, adminRole.ID); err != nil {
			return err
		}
	}

	return nil
}
