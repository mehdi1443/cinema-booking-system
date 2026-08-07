package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	cache_repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/random"
)

type authService struct {
	userRepo      repository_contract.UserInterface
	otpRepo       cache_repository_contract.OTPInterface
	refreshRepo   cache_repository_contract.RefreshTokenInterface
	tokenService  service_contract.JWTServiceInterface
	randomService *random.RandomGeneratorEngine
	otpConfig     *bootstrap.OTP
	jwtConfig     *bootstrap.JWT
	smsProvider   service_contract.SMSServiceProviderInterface
	emailProvider service_contract.EmailServiceProviderInterface
	hasher        *hash.BcryptHasher
	logger        logger.Logger
}

func NewAuthService(
	userRepo repository_contract.UserInterface,
	otpRepo cache_repository_contract.OTPInterface,
	refreshRepo cache_repository_contract.RefreshTokenInterface,
	tokenService service_contract.JWTServiceInterface,
	randomService *random.RandomGeneratorEngine,
	otpConfig *bootstrap.OTP,
	jwtConfig *bootstrap.JWT,
	smsProvider service_contract.SMSServiceProviderInterface,
	emailProvider service_contract.EmailServiceProviderInterface,
	hasher *hash.BcryptHasher,
	log logger.Logger,
) service_contract.AuthenticationServiceInterface {
	return &authService{
		userRepo:      userRepo,
		otpRepo:       otpRepo,
		refreshRepo:   refreshRepo,
		tokenService:  tokenService,
		randomService: randomService,
		otpConfig:     otpConfig,
		jwtConfig:     jwtConfig,
		smsProvider:   smsProvider,
		emailProvider: emailProvider,
		hasher:        hasher,
		logger:        log,
	}
}

func (s *authService) SendOTP(ctx context.Context, phone string, email string) (*bool, error) {
	if phone != "" {
		code := strconv.Itoa(s.randomService.GenerateRandomNumber(int(s.otpConfig.Digit)))

		err := s.otpRepo.SaveOTP(ctx, phone, code, s.otpConfig.ExpirationDuration())
		if err != nil {
			s.logger.Error(logger.Redis, logger.Insert, "failed to save phone OTP to cache", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to save OTP: %w", err)
		}

		s.smsProvider.SendSMS(dto.SendSMSRequest{
			Phone: phone,
			Text:  fmt.Sprintf("OTP:%s", code),
		})
		success := true
		return &success, nil
	} else if email != "" {
		code := strconv.Itoa(s.randomService.GenerateRandomNumber(int(s.otpConfig.Digit)))

		err := s.otpRepo.SaveOTP(ctx, email, code, s.otpConfig.ExpirationDuration())
		if err != nil {
			s.logger.Error(logger.Redis, logger.Insert, "failed to save email OTP to cache", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to save OTP: %w", err)
		}

		s.emailProvider.SendEmail(dto.SendEmailRequest{
			Email: email,
			Text:  fmt.Sprintf("OTP:%s", code),
		})
		success := true
		return &success, nil
	}
	return nil, fmt.Errorf("bad req to sendOTP service")
}

func (s *authService) VerifyOTP(ctx context.Context, phone string, code string) (*bool, error) {
	if phone == "" || code == "" {
		return nil, exception.ErrInvalidOTP
	}

	savedCode, err := s.otpRepo.GetOTP(ctx, phone)
	if err != nil {
		if errors.Is(err, exception.ErrOtpKeyNotFound) {
			s.logger.Warn(logger.Validation, logger.MobileValidation, "OTP verification attempted with expired code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, exception.ErrOTPExpired
		}
		return nil, err
	}

	if savedCode != code {
		s.logger.Warn(logger.Validation, logger.MobileValidation, "OTP verification failed due to code mismatch", nil)
		return nil, exception.ErrInvalidOTP
	}

	err = s.otpRepo.DeleteOTP(ctx, phone)
	if err != nil {
		s.logger.Error(logger.Redis, logger.Delete, "failed to delete OTP after verification", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	err = s.otpRepo.SetVerified(ctx, phone, s.otpConfig.VerifyExpirationDuration())
	if err != nil {
		s.logger.Error(logger.Redis, logger.Insert, "failed to mark phone as verified in cache", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	success := true
	return &success, nil
}

func (s *authService) Login(ctx context.Context, phone string, password string) (*dto.StartAuthResponse, error) {
	user, err := s.userRepo.FindUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			s.logger.Warn(logger.Internal, logger.Api, "login attempt failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, exception.ErrInvalidCredentials
		}
		return nil, err
	}

	// if user.IsPhoneVerified == false {
	// 	return nil, exception.ErrNotVerifiedPhoneNumber
	// }

	err = s.hasher.Compare(password, user.PasswordHash)
	if err != nil {
		s.logger.Warn(logger.Internal, logger.HashPassword, "login attempt failed: invalid password", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, exception.ErrInvalidCredentials
	}

	access, refresh, err := s.tokenService.GenerateTokens(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		s.logger.Error(logger.Internal, logger.Api, "failed to generate authentication tokens", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	strid := strconv.Itoa(int(user.ID))
	err = s.refreshRepo.SaveRefreshToken(ctx, strid, refresh, s.jwtConfig.RefreshDuration())
	if err != nil {
		s.logger.Error(logger.Redis, logger.Insert, "failed to store refresh token during login", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &dto.StartAuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
		Phone:        user.Phone,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}, nil
}

func (s *authService) Register(ctx context.Context, phone string, email string, password string) (*dto.StartAuthResponse, error) {
	existingUser, err := s.userRepo.FindUserByPhone(ctx, phone)
	if err == nil && existingUser != nil {
		return nil, exception.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return nil, err
	}

	existingUser, err = s.userRepo.FindUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, exception.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Error(logger.Internal, logger.HashPassword, "failed to hash password during registration", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &entity.User{
		Phone:           phone,
		Email:           email,
		FirstName:       "<Empty>",
		LastName:        "<Empty>",
		Status:          enum.UserStatusActive,
		WalletBalance:   0,
		IsPhoneVerified: false,
		IsEmailVerified: false,
		PasswordHash:    string(hashedPassword),
	}

	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Error(logger.Internal, logger.FailedToCreateUser, "failed to persist new user", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	access, refresh, err := s.tokenService.GenerateTokens(ctx, dto.GenerateTokenRequest{UserID: createdUser.ID})
	if err != nil {
		s.logger.Error(logger.Internal, logger.Api, "failed to generate tokens for newly registered user", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.StartAuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       newUser.ID,
		Phone:        newUser.Phone,
		Email:        newUser.Email,
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
	}, nil
}

func (s *authService) VerifyPhone(ctx context.Context, userID int64, phone string, code string) (*bool, error) {
	_, err := s.VerifyOTP(ctx, phone, code)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.VerifyUserPhone(ctx, userID)
	if err != nil {
		s.logger.Error(logger.Postgres, logger.Update, "failed to set user phone verified status", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	success := true
	return &success, nil
}

func (s *authService) VerifyEmail(ctx context.Context, tokenData *dto.TokenInContext, email string, code string) (*bool, error) {
	if email == "" || code == "" {
		return nil, exception.ErrInvalidOTP
	}

	savedCode, err := s.otpRepo.GetOTP(ctx, email)
	if err != nil {
		if errors.Is(err, exception.ErrOtpKeyNotFound) {
			s.logger.Warn(logger.Validation, logger.MobileValidation, "email verification attempted with expired code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, exception.ErrOTPExpired
		}
		return nil, err
	}

	if savedCode != code {
		s.logger.Warn(logger.Validation, logger.MobileValidation, "email verification code mismatch", nil)
		return nil, exception.ErrInvalidOTP
	}

	err = s.userRepo.VerifyUserEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.Postgres, logger.Update, "failed to set user email verified status", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	err = s.otpRepo.DeleteOTP(ctx, email)
	if err != nil {
		s.logger.Error(logger.Redis, logger.Delete, "failed to delete email OTP after verification", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, err
	}

	success := true
	return &success, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken *dto.RefreshTokenRequest) (*dto.StartAuthResponse, error) {
	if refreshToken.RefreshToken == "" {
		return nil, exception.ErrInvalidToken
	}

	claims, err := s.tokenService.ValidateToken(ctx, refreshToken.RefreshToken, string(domainjwt.RefreshTokenTag))
	if err != nil {
		s.logger.Warn(logger.Internal, logger.Api, "refresh token validation failed", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, exception.ErrTokenExpired
	}

	storedToken, err := s.refreshRepo.GetRefreshToken(ctx, strconv.Itoa(int(claims.UserID)))
	if err != nil || storedToken != refreshToken.RefreshToken {
		s.logger.Warn(logger.Internal, logger.Api, "refresh token mismatch or missing in store", nil)
		return nil, exception.ErrInvalidToken
	}

	user, err := s.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			return nil, exception.ErrUserNotFound
		}
		return nil, err
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		s.logger.Error(logger.Internal, logger.Api, "failed to generate new access token during refresh", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		s.logger.Error(logger.Internal, logger.Api, "failed to generate new refresh token during refresh", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.refreshRepo.SaveRefreshToken(ctx, strconv.Itoa(int(user.ID)), newRefreshToken, s.jwtConfig.RefreshDuration())
	if err != nil {
		s.logger.Error(logger.Redis, logger.Insert, "failed to store updated refresh token", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.StartAuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       user.ID,
		Phone:        user.Phone,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID int64, jwtContext *dto.JWTTokenInContext, refreshToken *dto.LogoutRequest) error {
	if userID == 0 {
		return exception.ErrInvalidCredentials
	}

	err := s.refreshRepo.DeleteRefreshToken(ctx, strconv.Itoa(int(userID)))
	if err != nil {
		s.logger.Error(logger.Redis, logger.Delete, "failed to remove refresh token on logout", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	err = s.tokenService.BlacklistToken(ctx, jwtContext.TheToken, jwtContext.TheTokenExpireAt.Unix())
	if err != nil {
		s.logger.Error(logger.Redis, logger.Insert, "failed to blacklist access token on logout", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return err
	}

	return nil
}
