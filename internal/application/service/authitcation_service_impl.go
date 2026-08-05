package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"github.com/mehdi1443/cinema-booking-system/pkg/random"
)

type authService struct {
	userRepo      repository_contract.UserInterFace
	otpRepo       cache_repository_contract.OTPInterface
	refreshRepo   cache_repository_contract.RefreshTokenInterface
	tokenService  service_contract.JWTServiceInterface
	randomService *random.RandomGeneratorEngine
	otpConfig     *bootstrap.OTP
	jwtConfig     *bootstrap.JWT
	smsProvider   service_contract.SMSServiceProviderInterface
	hasher        *hash.BcryptHasher
}

func NewAuthService(
	userRepo repository_contract.UserInterFace,
	otpRepo cache_repository_contract.OTPInterface,
	refreshRepo cache_repository_contract.RefreshTokenInterface,
	tokenService service_contract.JWTServiceInterface,
	randomService *random.RandomGeneratorEngine,
	otpConfig *bootstrap.OTP,
	jwtConfig *bootstrap.JWT,
	smsProvider service_contract.SMSServiceProviderInterface,
	hasher *hash.BcryptHasher,
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
		hasher:        hasher,
	}
}

func (s *authService) SendOTP(ctx context.Context, phone string) (*bool, error) {
	if phone == "" {
		return nil, exception.ErrInvalidPhone
	}

	code := strconv.Itoa(s.randomService.GenerateRandomNumber(int(s.otpConfig.Digit)))

	err := s.otpRepo.SaveOTP(ctx, phone, code, s.otpConfig.ExpirationDuration())
	if err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}

	//TODO: Send OTP via SMS provider logic here...
	s.smsProvider.SendSMS(dto.SendSMSRequest{
		Phone: phone,
		Text:  fmt.Sprintf("OTP:%s", code),
	})
	success := true
	return &success, nil
}

func (s *authService) VerifyOTP(ctx context.Context, phone string, code string) (*bool, error) {
	if phone == "" || code == "" {
		return nil, exception.ErrInvalidOTP
	}

	savedCode, err := s.otpRepo.GetOTP(ctx, phone)
	if err != nil {
		if errors.Is(err, exception.ErrOtpKeyNotFound) {
			return nil, exception.ErrOTPExpired
		}
		return nil, err
	}

	if savedCode != code {
		return nil, exception.ErrInvalidOTP
	}

	// Invalidate OTP after successful verification
	err = s.otpRepo.DeleteOTP(ctx, phone)
	if err != nil {
		return nil, err
	}
	err = s.otpRepo.SetVerified(ctx, phone, s.otpConfig.VerifyExpirationDuration())
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}

func (s *authService) Login(ctx context.Context, phone string, password string) (access string, refresh string, err error) {
	user, err := s.userRepo.FindUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			log.Printf("loginsrv:%v\n", err)
			return "", "", exception.ErrInvalidCredentials
		}
		return "", "", err
	}
	if user.IsPhoneVerified == false {
		return "", "", exception.ErrNotVerifiedPhoneNumber
	}

	err = s.hasher.Compare(password, user.PasswordHash)
	if err != nil {
		log.Printf("loginsrv:%v\n", err)
		return "", "", exception.ErrInvalidCredentials
	}

	access, refresh, err = s.tokenService.GenerateTokens(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}
	strid := strconv.Itoa(int(user.ID))

	err = s.refreshRepo.SaveRefreshToken(ctx, strid, refresh, s.jwtConfig.RefreshDuration())
	if err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}
	return access, refresh, nil
}

func (s *authService) Register(ctx context.Context, phone string, email string, password string) (access string, refresh string, err error) {
	existingUser, err := s.userRepo.FindUserByPhone(ctx, phone)
	if err == nil && existingUser != nil {
		return "", "", exception.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return "", "", err
	}

	existingUser, err = s.userRepo.FindUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return "", "", exception.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return "", "", err
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
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
		return "", "", err
	}

	access, refresh, err = s.tokenService.GenerateTokens(ctx, dto.GenerateTokenRequest{UserID: createdUser.ID})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return access, refresh, nil
}

func (s *authService) VerifyPhone(ctx context.Context, tokenData *dto.TokenInContext, phone string, code string) (*bool, error) {
	_, err := s.VerifyOTP(ctx, phone, code)
	if err != nil {
		return nil, err
	}
	err = s.userRepo.VerifyUserPhone(ctx, tokenData.UserID)
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}

// TODO : send otp email handler/srv
func (s *authService) VerifyEmail(ctx context.Context, tokenData *dto.TokenInContext, email string, code string) (*bool, error) {
	if email == "" || code == "" {
		return nil, exception.ErrInvalidOTP
	}

	savedCode, err := s.otpRepo.GetOTP(ctx, email)
	if err != nil {
		if errors.Is(err, exception.ErrOtpKeyNotFound) {
			return nil, exception.ErrOTPExpired
		}
		return nil, err
	}

	if savedCode != code {
		return nil, exception.ErrInvalidOTP
	}

	err = s.userRepo.VerifyUserEmail(ctx, tokenData.UserID)
	if err != nil {
		return nil, err
	}

	err = s.otpRepo.DeleteOTP(ctx, email)
	if err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}

// RefreshToken validates the refresh token and generates a new token pair.
func (s *authService) RefreshToken(ctx context.Context, refreshToken *dto.RefreshTokenRequest) (*dto.AuthTokenResponse, error) {
	if refreshToken.RefreshToken == "" {
		return nil, exception.ErrInvalidToken
	}

	claims, err := s.tokenService.ValidateToken(ctx, refreshToken.RefreshToken, string(domainjwt.RefreshTokenTag))
	if err != nil {
		return nil, exception.ErrTokenExpired
	}

	storedToken, err := s.refreshRepo.GetRefreshToken(ctx, strconv.Itoa(int(claims.UserID)))
	if err != nil || storedToken != refreshToken.RefreshToken {
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
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.refreshRepo.SaveRefreshToken(ctx, strconv.Itoa(int(user.ID)), newRefreshToken, s.jwtConfig.RefreshDuration())
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.AuthTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout invalidates the refresh token and terminates the session.
func (s *authService) Logout(ctx context.Context, tokenData *dto.TokenInContext, jwtContext *dto.JWTTokenInContext, refreshToken *dto.LogoutRequest) error {
	if tokenData == nil || tokenData.UserID == 0 {
		return exception.ErrInvalidCredentials
	}

	err := s.refreshRepo.DeleteRefreshToken(ctx, strconv.Itoa(int(tokenData.UserID)))
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	err = s.tokenService.BlacklistToken(ctx, jwtContext.TheToken, jwtContext.TheTokenExpireAt.Unix())
	if err != nil {
		return err
	}
	return nil
}
