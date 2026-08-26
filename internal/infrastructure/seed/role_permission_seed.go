package seed

import (
	"context"
	"errors"
	"fmt"
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

		for _, permType := range roleEnum.DefaultPermissions() {
			if permType == enum.PermissionAll {
				for _, perm := range permMap {
					if err := s.appendPermissionIfNotExist(ctx, role.ID, perm.ID); err != nil {
						return nil, fmt.Errorf("append perm to role: %w", err)
					}
				}
			} else if perm, exists := permMap[permType]; exists {
				if err := s.appendPermissionIfNotExist(ctx, role.ID, perm.ID); err != nil {
					return nil, fmt.Errorf("append perm to role: %w", err)
				}
			}
		}
	}

	return roleMap, nil
}

func (s *RoleSeeder) appendPermissionIfNotExist(ctx context.Context, roleID, permID int64) error {
	hasPerm, err := s.rbacRepo.RoleHasPermission(ctx, roleID, permID)
	if err != nil {
		return err
	}
	if !hasPerm {
		return s.rbacRepo.AppendPermissionToRole(ctx, roleID, permID)
	}
	return nil
}

func (s *RoleSeeder) seedSuperAdmin(ctx context.Context, adminRole *entity.Role) error {
	user, err := s.userRepo.FindUserByPhone(ctx, s.superAdmin.Phone)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return err
	}

	if errors.Is(err, exception.ErrUserNotFound) {
		hashedPassword, err := s.hasher.Hash(s.superAdmin.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user = &entity.User{
			Phone:           s.superAdmin.Phone,
			Email:           s.superAdmin.Email,
			PasswordHash:    hashedPassword,
			FirstName:       s.superAdmin.FirstName,
			LastName:        s.superAdmin.LastName,
			WalletBalance:   s.superAdmin.Wallet,
			IsPhoneVerified: true,
			IsEmailVerified: true,
			PhoneVerifiedAt: time.Now(),
			EmailVerifiedAt: time.Now(),
			Status:          "ACTIVE",
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
