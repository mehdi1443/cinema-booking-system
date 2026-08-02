package jwtsrv

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/pkg/rsa"
)

type jwtService struct {
	rsaKeys *rsa.RSAKeys
	cfg     *bootstrap.JWT
}

func NewTokenService(rsaKeys *rsa.RSAKeys, cfg *bootstrap.JWT) service_contract.JWTServiceInterface {
	return &jwtService{
		rsaKeys: rsaKeys,
		cfg:     cfg,
	}
}

func (s *jwtService) GenerateAccessToken(req dto.GenerateTokenRequest) (string, error) {
	now := time.Now()
	payload := domainjwt.TokenPayload{
		UserID:    req.UserID,
		TokenType: domainjwt.AccessTokenTag,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.cfg.AccessDuration()),
	}
	return s.generateToken(payload)
}

func (s *jwtService) GenerateRefreshToken(req dto.GenerateTokenRequest) (string, error) {
	now := time.Now()
	payload := domainjwt.TokenPayload{
		UserID:    req.UserID,
		TokenType: domainjwt.RefreshTokenTag,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.cfg.RefreshDuration()),
	}
	return s.generateToken(payload)
}

func (s *jwtService) ValidateToken(signedToken string) (*domainjwt.TokenPayload, error) {
	parsedToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.rsaKeys.GetPublicKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, service_contract.ErrTokenExpired
		}
		return nil, errors.Join(service_contract.ErrInvalidToken, err)
	}

	if parsedToken == nil || !parsedToken.Valid {
		return nil, service_contract.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, service_contract.ErrInvalidToken
	}

	return s.parseClaimsToPayload(claims)
}

// Private helper to sign claims with RS256
func (s *jwtService) generateToken(payload domainjwt.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		service_contract.UserIDKeyMapClaims:    payload.UserID,
		service_contract.TokenTypeKeyMapClaims: string(payload.TokenType),
		service_contract.IssuedAtKeyMapClaims:  payload.IssuedAt.Unix(),
		service_contract.ExpireAtKeyMapClaims:  payload.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.rsaKeys.GetPrivateKey())
	if err != nil {
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
