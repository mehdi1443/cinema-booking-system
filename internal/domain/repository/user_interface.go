package repository_contract

import (
	"context"
	"errors"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

var (
	// User Errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Role Errors
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")

	// Permission Errors
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")

	// Generic Storage Error
	ErrDatabase = errors.New("database operation failed")
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

	CreatePermission(ctx context.Context, permission *entity.Permission) error
	FindPermissionByName(ctx context.Context, name string) (*entity.Permission, error)

	AssignPermissionToRole(ctx context.Context, roleID int64, permissionID int64) error
	AssignRoleToUser(ctx context.Context, userID int64, roleID int64) error
	RoleHasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error)
	UserHasRole(ctx context.Context, userID int64, roleID int64) (bool, error)

	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}
