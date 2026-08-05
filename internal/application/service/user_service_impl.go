package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
)

type userService struct {
	userRepo repository_contract.UserInterFace
	hasher   *hash.BcryptHasher
}

func NewUserService(userRepo repository_contract.UserInterFace) service_contract.UserServiceInterface {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, ctxToken *dto.TokenInContext, id int64, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user, hideOpt), nil
}

func (s *userService) GetUserByPhone(ctx context.Context, ctxToken *dto.TokenInContext, phone string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user, hideOpt), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, ctxToken *dto.TokenInContext, email string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user, hideOpt), nil
}

func (s *userService) GetUsersByStatus(ctx context.Context, ctxToken *dto.TokenInContext, status enum.UserStatus, hideOpt *dto.HideFieldUserOption) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.FindUsersByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := dto.ToUserResponseSlice(users, hideOpt)
	return responses, nil
}

func (s *userService) UpdateProfile(ctx context.Context, ctxToken *dto.TokenInContext, req dto.UpdateUserProfileDTO) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email

	updatedUser, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(updatedUser, &dto.HideFieldUserOption{WalletBalance: false}), nil
}



func (s *userService) ChangeStatus(ctx context.Context, ctxToken *dto.TokenInContext, userID int64, status enum.UserStatus) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Status = status

	updatedUser, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(updatedUser, &dto.HideFieldUserOption{WalletBalance: false}), nil
}

func (s *userService) DeleteUserByID(ctx context.Context, ctxToken *dto.TokenInContext, id int64) error {
	// Verify record exists before deletion
	_, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepo.DeleteUserByID(ctx, id)
}

func (s *userService) GetUsers(ctx context.Context, tokenData *dto.TokenInContext, query *dto.GetUsersQueryDTO) ([]*dto.UserResponse, error) {
	if query == nil || query.IsEmpty() {
		return nil, exception.ErrInvalidFilter
	}

	users, err := s.userRepo.FindUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, exception.ErrUserNotFound
	}

	return dto.ToUserResponseSlice(users, &query.HideFieldUserOption), nil
}

func (s *userService) UpdatePassword(ctx context.Context, ctxToken *dto.TokenInContext, userId int64, password dto.UpdatePasswordRequest) (*bool, error) {
	if password.Password == "" {
		return nil, exception.ErrInvalidPassword
	}

	// 1. Verify user existence
	_, err := s.userRepo.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 2. Hash the raw password
	hashedPassword, err := s.hasher.Hash(password.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Update the password in the database
	err = s.userRepo.UpdatePassword(ctx, userId, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}
