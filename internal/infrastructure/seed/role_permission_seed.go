package seed

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
	"gorm.io/gorm"
)

type RoleSeeder struct {
	superAdmin bootstrap.SuperUser
	userRepo   repository_contract.UserInterface
	rbacRepo   repository_contract.RBACInterface
	db         database.Database
	hasher     *hash.BcryptHasher
}

func NewRoleSeeder(
	superAdmin bootstrap.SuperUser,
	userRepo repository_contract.UserInterface,
	rbacRepo repository_contract.RBACInterface,
	db database.Database,
	hasher *hash.BcryptHasher,
) *RoleSeeder {
	return &RoleSeeder{
		superAdmin: superAdmin,
		userRepo:   userRepo,
		rbacRepo:   rbacRepo,
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
		perm, err := s.rbacRepo.FindPermissionByName(ctx, permName)
		if err != nil && !errors.Is(err, exception.ErrPermissionNotFound) {
			return nil, fmt.Errorf("find permission %s: %w", permName, err)
		}

		if errors.Is(err, exception.ErrPermissionNotFound) {
			perm = &entity.Permission{
				Name:        permName,
				Title:       permType.Description(),
				Category:    permType.Category(),
				Description: permType.Description(),
			}
			if err := s.rbacRepo.CreatePermission(ctx, perm); err != nil {
				return nil, fmt.Errorf("create permission %s: %w", permName, err)
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
		role, err := s.rbacRepo.FindRoleByName(ctx, roleEnum.String())
		if err != nil && !errors.Is(err, exception.ErrRoleNotFound) {
			return nil, fmt.Errorf("find role %s: %w", roleEnum.String(), err)
		}

		if errors.Is(err, exception.ErrRoleNotFound) {
			role = &entity.Role{
				Name:  roleEnum.String(),
				Title: roleEnum.Title(),
			}
			if err := s.rbacRepo.CreateRole(ctx, role); err != nil {
				return nil, fmt.Errorf("create role %s: %w", roleEnum.String(), err)
			}
		}
		roleMap[roleEnum] = role

		permissionIDs := make([]int64, 0, len(roleEnum.DefaultPermissions()))
		for _, permType := range roleEnum.DefaultPermissions() {
			if perm, exists := permMap[permType]; exists {
				permissionIDs = append(permissionIDs, perm.ID)
			}
		}
		if err := s.rbacRepo.ReplacePermissionsOfRole(ctx, role.ID, permissionIDs); err != nil {
			return nil, fmt.Errorf("synchronize permissions for role %s: %w", role.Name, err)
		}
	}

	return roleMap, nil
}

func (s *RoleSeeder) seedSuperAdmin(ctx context.Context, adminRole *entity.Role) error {
	if strings.TrimSpace(s.superAdmin.Phone) == "" || strings.TrimSpace(s.superAdmin.Password) == "" {
		admins, err := s.rbacRepo.GetUsersByRoleNames(ctx, []string{enum.RoleAdmin.String()})
		if err != nil {
			return err
		}
		for _, admin := range admins {
			if admin != nil && admin.Status == enum.UserActive {
				return nil
			}
		}
		return errors.New("SUPER_ADMIN_PHONE and SUPER_ADMIN_PASSWORD must be configured when no active administrator exists")
	}
	user, err := s.userRepo.FindUserByPhone(ctx, s.superAdmin.Phone)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return err
	}

	if errors.Is(err, exception.ErrUserNotFound) {
		hashedPassword, err := s.hasher.Hash(s.superAdmin.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		now := time.Now()

		var emailPtr *string
		if s.superAdmin.Email != "" {
			email := s.superAdmin.Email
			emailPtr = &email
		}

		user = &entity.User{
			Phone:           s.superAdmin.Phone,
			Email:           emailPtr,
			PasswordHash:    &hashedPassword,
			FirstName:       s.superAdmin.FirstName,
			LastName:        s.superAdmin.LastName,
			WalletBalance:   s.superAdmin.Wallet,
			IsPhoneVerified: true,
			IsEmailVerified: true,
			PhoneVerifiedAt: &now,
			EmailVerifiedAt: &now,
			Status:          enum.UserActive,
		}

		user, err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			if errors.Is(err, exception.ErrUserAlreadyExists) || errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil
			}
			return err
		}
	}

	hasRole, err := s.rbacRepo.UserHasRole(ctx, user.ID, adminRole.ID)
	if err != nil {
		return err
	}
	if !hasRole {
		if err := s.rbacRepo.AppendRoleToUser(ctx, user.ID, adminRole.ID); err != nil {
			return err
		}
	}

	return nil
}
