package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type UserServiceInterface interface {
	GetUserByID(ctx context.Context, ctxToken *dto.TokenInContext, id int64, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)
	GetUserByPhone(ctx context.Context, ctxToken *dto.TokenInContext, phone string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, ctxToken *dto.TokenInContext, email string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)

	GetUsersByStatus(ctx context.Context, ctxToken *dto.TokenInContext, status enum.UserStatus, hideOpt *dto.HideFieldUserOption) ([]*dto.UserResponse, error)

	UpdateProfile(ctx context.Context, ctxToken *dto.TokenInContext, userID int64, dto dto.UpdateUserProfileDTO) (*dto.UserResponse, error)
	ChangeStatus(ctx context.Context, ctxToken *dto.TokenInContext, userID int64, status enum.UserStatus) (*dto.UserResponse, error)

	DeleteUserByID(ctx context.Context, ctxToken *dto.TokenInContext, id int64) error
	GetUsers(ctx context.Context, tokenData *dto.TokenInContext, query *dto.GetUsersQueryDTO) ([]*dto.UserResponse, error)

	UpdatePassword(ctx context.Context, ctxToken *dto.TokenInContext, userId int64, password dto.UpdatePasswordRequest) (*bool, error)
}
