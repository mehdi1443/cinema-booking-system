package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/rsa"
)

type jwtService struct {
	accessRSAKeys  *rsa.RSAKeys
	refreshRSAKeys *rsa.RSAKeys
	cfg            *bootstrap.JWT
	blacklistRepo  cache.TokenBlacklistInterface
	logger         domainlogger.Logger
}

func NewTokenService(
	accessRSAKeys *rsa.RSAKeys,
	refreshRSAKeys *rsa.RSAKeys,
	cfg *bootstrap.JWT,
	blacklistRepo cache.TokenBlacklistInterface,
	log domainlogger.Logger,
) service_contract.JWTServiceInterface {
	return &jwtService{
		accessRSAKeys:  accessRSAKeys,
		refreshRSAKeys: refreshRSAKeys,
		cfg:            cfg,
		blacklistRepo:  blacklistRepo,
		logger:         log,
	}
}

func (s *jwtService) GenerateTokens(ctx context.Context, req dto.GenerateTokenRequest) (access string, refresh string, err error) {
	access, err = s.GenerateAccessToken(ctx, req)
	if err != nil {
		return "", "", err
	}
	refresh, err = s.GenerateRefreshToken(ctx, req)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *jwtService) GenerateAccessToken(ctx context.Context, req dto.GenerateTokenRequest) (string, error) {
	now := time.Now()
	payload := domainjwt.TokenPayload{
		UserID:    req.UserID,
		TokenType: domainjwt.AccessTokenTag,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.cfg.AccessDuration()),
	}
	return s.generateToken(ctx, payload)
}

func (s *jwtService) GenerateRefreshToken(ctx context.Context, req dto.GenerateTokenRequest) (string, error) {
	now := time.Now()
	payload := domainjwt.TokenPayload{
		UserID:    req.UserID,
		TokenType: domainjwt.RefreshTokenTag,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.cfg.RefreshDuration()),
	}
	return s.generateToken(ctx, payload)
}

func (s *jwtService) ValidateToken(ctx context.Context, signedToken string, tokenType string) (*domainjwt.TokenPayload, error) {
	requestedType := domainjwt.TokenType(tokenType)
	if requestedType != domainjwt.AccessTokenTag && requestedType != domainjwt.RefreshTokenTag {
		return nil, exception.ErrInvalidToken
	}

	parsedToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		if requestedType == domainjwt.AccessTokenTag {
			return s.accessRSAKeys.GetPublicKey(), nil
		}

		return s.refreshRSAKeys.GetPublicKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, exception.ErrTokenExpired
		}
		s.logger.Warn(ctx, "failed to parse JWT token", domainlogger.Error(err))
		return nil, errors.Join(exception.ErrInvalidToken, err)
	}

	if parsedToken == nil || !parsedToken.Valid {
		return nil, exception.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, exception.ErrInvalidToken
	}

	payload, err := s.parseClaimsToPayload(claims)
	if err != nil {
		return nil, err
	}
	if payload.TokenType != requestedType {
		return nil, exception.ErrInvalidToken
	}
	return payload, nil
}

func (s *jwtService) generateToken(ctx context.Context, payload domainjwt.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		service_contract.UserIDKeyMapClaims:    payload.UserID,
		service_contract.TokenTypeKeyMapClaims: string(payload.TokenType),
		service_contract.IssuedAtKeyMapClaims:  payload.IssuedAt.Unix(),
		service_contract.ExpireAtKeyMapClaims:  payload.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	var signedToken string
	var err error
	if payload.TokenType == domainjwt.AccessTokenTag {
		signedToken, err = token.SignedString(s.accessRSAKeys.GetPrivateKey())
	} else if payload.TokenType == domainjwt.RefreshTokenTag {
		signedToken, err = token.SignedString(s.refreshRSAKeys.GetPrivateKey())
	} else {
		s.logger.Error(ctx, "unknown token type provided for signing", domainlogger.String("type", string(payload.TokenType)))
		return "", fmt.Errorf("unsupported token type %q", payload.TokenType)
	}

	if err != nil {
		s.logger.Error(ctx, "failed to sign JWT string with private RSA key", domainlogger.Error(err))
		return "", fmt.Errorf("jwtsrv.generateToken: failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *jwtService) parseClaimsToPayload(claims jwt.MapClaims) (*domainjwt.TokenPayload, error) {
	payload := &domainjwt.TokenPayload{}

	if uid, ok := claims[service_contract.UserIDKeyMapClaims].(float64); ok {
		payload.UserID = int64(uid)
	} else if uidInt, ok := claims[service_contract.UserIDKeyMapClaims].(int64); ok {
		payload.UserID = uidInt
	}

	if tType, ok := claims[service_contract.TokenTypeKeyMapClaims].(string); ok {
		payload.TokenType = domainjwt.TokenType(tType)
	}

	if exp, ok := claims[service_contract.ExpireAtKeyMapClaims].(float64); ok {
		payload.ExpiresAt = time.Unix(int64(exp), 0)
	}

	if iat, ok := claims[service_contract.IssuedAtKeyMapClaims].(float64); ok {
		payload.IssuedAt = time.Unix(int64(iat), 0)
	}

	return payload, nil
}

func (s *jwtService) BlacklistToken(ctx context.Context, token string, expiresAt int64) error {
	if token == "" {
		return nil
	}

	expirationTime := time.Unix(expiresAt, 0)
	ttl := time.Until(expirationTime)

	if ttl <= 0 {
		return nil
	}

	err := s.blacklistRepo.BlacklistToken(ctx, token, ttl)
	if err != nil {
		s.logger.Error(ctx, "failed to record token in blacklist cache", domainlogger.Error(err))
		return err
	}
	return nil
}

func (s *jwtService) IsBlacklisted(ctx context.Context, token string) (*bool, error) {
	if token == "" {
		success := false
		return &success, nil
	}
	return s.blacklistRepo.IsBlacklisted(ctx, token)
}
