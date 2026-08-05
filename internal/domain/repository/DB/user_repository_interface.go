package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type UserInterFace interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByID(ctx context.Context, id int64) (*entity.User, error)
	FindUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)

	FindUsersByStatus(ctx context.Context, status enum.UserStatus) ([]*entity.User, error)
	DeleteUserByID(ctx context.Context, id int64) error
	FindUsers(ctx context.Context, filter *dto.GetUsersQueryDTO) ([]*entity.User, error)

	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error

	VerifyUserPhone(ctx context.Context, userID int64) error
	VerifyUserEmail(ctx context.Context, userID int64) error
}
