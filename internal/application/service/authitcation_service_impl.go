package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	cache_repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
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
	debugConfig   *bootstrap.DebugModeOptions
	smsProvider   service_contract.SMSServiceProviderInterface
	emailProvider service_contract.EmailServiceProviderInterface
	hasher        *hash.BcryptHasher
	logger        domainlogger.Logger
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
	log domainlogger.Logger,
) service_contract.AuthenticationServiceInterface {
	return newAuthService(userRepo, otpRepo, refreshRepo, tokenService, randomService, otpConfig, jwtConfig, nil, smsProvider, emailProvider, hasher, log)
}

func NewAuthServiceWithDebug(
	userRepo repository_contract.UserInterface,
	otpRepo cache_repository_contract.OTPInterface,
	refreshRepo cache_repository_contract.RefreshTokenInterface,
	tokenService service_contract.JWTServiceInterface,
	randomService *random.RandomGeneratorEngine,
	otpConfig *bootstrap.OTP,
	jwtConfig *bootstrap.JWT,
	debugConfig *bootstrap.DebugModeOptions,
	smsProvider service_contract.SMSServiceProviderInterface,
	emailProvider service_contract.EmailServiceProviderInterface,
	hasher *hash.BcryptHasher,
	log domainlogger.Logger,
) service_contract.AuthenticationServiceInterface {
	return newAuthService(userRepo, otpRepo, refreshRepo, tokenService, randomService, otpConfig, jwtConfig, debugConfig, smsProvider, emailProvider, hasher, log)
}

func newAuthService(
	userRepo repository_contract.UserInterface,
	otpRepo cache_repository_contract.OTPInterface,
	refreshRepo cache_repository_contract.RefreshTokenInterface,
	tokenService service_contract.JWTServiceInterface,
	randomService *random.RandomGeneratorEngine,
	otpConfig *bootstrap.OTP,
	jwtConfig *bootstrap.JWT,
	debugConfig *bootstrap.DebugModeOptions,
	smsProvider service_contract.SMSServiceProviderInterface,
	emailProvider service_contract.EmailServiceProviderInterface,
	hasher *hash.BcryptHasher,
	log domainlogger.Logger,
) service_contract.AuthenticationServiceInterface {
	return &authService{
		userRepo: userRepo, otpRepo: otpRepo, refreshRepo: refreshRepo,
		tokenService: tokenService, randomService: randomService,
		otpConfig: otpConfig, jwtConfig: jwtConfig, debugConfig: debugConfig, smsProvider: smsProvider,
		emailProvider: emailProvider, hasher: hasher, logger: log,
	}
}

func normalizeIranianPhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	switch {
	case strings.HasPrefix(phone, "+98"):
		phone = "0" + strings.TrimPrefix(phone, "+98")
	case len(phone) == 10 && strings.HasPrefix(phone, "9"):
		phone = "0" + phone
	}
	if len(phone) != 11 || !strings.HasPrefix(phone, "09") {
		return "", exception.ErrInvalidPhone
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return "", exception.ErrInvalidPhone
		}
	}
	return phone, nil
}

func (s *authService) SendOTP(ctx context.Context, rawPhone string) (*dto.SendOTPResponse, error) {
	phone, err := normalizeIranianPhone(rawPhone)
	if err != nil {
		return nil, err
	}
	code := strconv.Itoa(s.randomService.GenerateRandomNumber(int(s.otpConfig.Digit)))
	if err := s.otpRepo.SaveOTP(ctx, phone, code, s.otpConfig.ExpirationDuration()); err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}
	if _, err := s.smsProvider.SendSMS(dto.SendSMSRequest{Phone: phone, Text: fmt.Sprintf("OTP:%s", code)}); err != nil {
		_ = s.otpRepo.DeleteOTP(ctx, phone)
		return nil, fmt.Errorf("failed to send OTP: %w", err)
	}
	response := &dto.SendOTPResponse{
		Phone:            phone,
		ExpiresInSeconds: int64(s.otpConfig.ExpirationDuration().Seconds()),
	}
	// The OTP is returned only for the local development flow. A real SMS
	// provider must never expose the authentication secret in the API response.
	if s.debugConfig != nil && s.debugConfig.Enable {
		response.OTPCode = code
	}
	return response, nil
}

func (s *authService) VerifyOTP(ctx context.Context, rawPhone string, code string) (*dto.StartAuthResponse, error) {
	phone, err := normalizeIranianPhone(rawPhone)
	if err != nil {
		return nil, err
	}
	if code == "" {
		return nil, exception.ErrInvalidOTP
	}
	if err := s.otpRepo.ConsumeOTP(ctx, phone, code); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByPhone(ctx, phone)
	isNewUser := false
	if errors.Is(err, exception.ErrUserNotFound) {
		now := time.Now()
		user = &entity.User{
			Phone: phone, Status: enum.UserActive, WalletBalance: 0,
			IsPhoneVerified: true, PhoneVerifiedAt: &now,
		}
		user, err = s.userRepo.CreateUserWithRole(ctx, user, enum.RoleUser.String())
		if err != nil {
			return nil, fmt.Errorf("create OTP user: %w", err)
		}
		isNewUser = true
	} else if err != nil {
		return nil, err
	} else {
		if user.Status != enum.UserActive {
			return nil, exception.ErrUnauthorized
		}
		if !user.IsPhoneVerified {
			if err := s.userRepo.VerifyUserPhone(ctx, user.ID); err != nil {
				return nil, err
			}
			user.IsPhoneVerified = true
		}
	}
	return s.issueSession(ctx, user, isNewUser)
}

func (s *authService) issueSession(ctx context.Context, user *entity.User, isNewUser bool) (*dto.StartAuthResponse, error) {
	access, refresh, err := s.tokenService.GenerateTokens(ctx, dto.GenerateTokenRequest{UserID: user.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	if err := s.refreshRepo.SaveRefreshToken(ctx, strconv.FormatInt(user.ID, 10), refresh, s.jwtConfig.RefreshDuration()); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		if role != nil {
			roles = append(roles, role.Name)
		}
	}
	return &dto.StartAuthResponse{
		AccessToken: access, RefreshToken: refresh, UserID: user.ID,
		Phone: user.Phone, Email: user.Email, FirstName: user.FirstName,
		LastName: user.LastName, IsNewUser: isNewUser, Roles: roles,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, request *dto.RefreshTokenRequest) (*dto.StartAuthResponse, error) {
	if request == nil || request.RefreshToken == "" {
		return nil, exception.ErrInvalidToken
	}
	claims, err := s.tokenService.ValidateToken(ctx, request.RefreshToken, string(domainjwt.RefreshTokenTag))
	if err != nil {
		if errors.Is(err, exception.ErrTokenExpired) {
			return nil, exception.ErrTokenExpired
		}
		return nil, exception.ErrInvalidToken
	}
	storedToken, err := s.refreshRepo.GetRefreshToken(ctx, strconv.FormatInt(claims.UserID, 10))
	if err != nil || storedToken != request.RefreshToken {
		return nil, exception.ErrInvalidToken
	}
	user, err := s.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status != enum.UserActive || !user.IsPhoneVerified {
		return nil, exception.ErrUnauthorized
	}
	return s.issueSession(ctx, user, false)
}

func (s *authService) Logout(ctx context.Context, userID int64, jwtContext *dto.JWTTokenInContext, _ *dto.LogoutRequest) error {
	if userID <= 0 || jwtContext == nil {
		return exception.ErrInvalidCredentials
	}
	if err := s.refreshRepo.DeleteRefreshToken(ctx, strconv.FormatInt(userID, 10)); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return s.tokenService.BlacklistToken(ctx, jwtContext.TheToken, jwtContext.TheTokenExpireAt.Unix())
}
