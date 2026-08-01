package repository

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
)

type UserInterFace interface {
	CreateUser(ctx context.Context, db database.Database, user *entity.User) (*entity.User, error)
	FindUserByID(ctx context.Context, db database.Database, id int64) (*entity.User, error)
	FindUserByPhone(ctx context.Context, db database.Database, phone string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, db database.Database, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, db database.Database, user *entity.User) (*entity.User, error)

	CreateRole(ctx context.Context, db database.Database, role *entity.Role) error
	FindRoleByName(ctx context.Context, db database.Database, name string) (*entity.Role, error)
	FindRoleByID(ctx context.Context, db database.Database, roleID int64) (*entity.Role, error)

	CreatePermission(ctx context.Context, db database.Database, permission *entity.Permission) error
	FindPermissionByName(ctx context.Context, db database.Database, name string) (*entity.Permission, error)

	AssignPermissionToRole(ctx context.Context, db database.Database, roleID int64, permissionID int64) error
	AssignRoleToUser(ctx context.Context, db database.Database, userID int64, roleID int64) error
	RoleHasPermission(ctx context.Context, db database.Database, roleID int64, permissionID int64) (bool, error)
	UserHasRole(ctx context.Context, db database.Database, userID int64, roleID int64) (bool, error)

	GetUserPermissions(ctx context.Context, db database.Database, userID int64) ([]string, error)
}
