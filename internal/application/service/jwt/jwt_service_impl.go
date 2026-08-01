package jwtsrv

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application"
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
		Roles:     req.Roles,
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
		Roles:     req.Roles,
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

	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return s.parseClaimsToPayload(claims)
}

// Private helper to sign claims with RS256
func (s *jwtService) generateToken(payload domainjwt.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    payload.UserID,
		"roles":      payload.Roles,
		"token_type": string(payload.TokenType),
		"iat":        payload.IssuedAt.Unix(),
		"exp":        payload.ExpiresAt.Unix(),
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

	if uid, ok := claims["user_id"].(float64); ok {
		payload.UserID = int64(uid)
	} else if uidInt, ok := claims["user_id"].(int64); ok {
		payload.UserID = uidInt
	}

	if rawRoles, ok := claims["roles"].([]interface{}); ok {
		roles := make([]string, len(rawRoles))
		for i, r := range rawRoles {
			if str, isStr := r.(string); isStr {
				roles[i] = str
			}
		}
		payload.Roles = roles
	}

	if tType, ok := claims["token_type"].(string); ok {
		payload.TokenType = domainjwt.TokenType(tType)
	}

	if exp, ok := claims["exp"].(float64); ok {
		payload.ExpiresAt = time.Unix(int64(exp), 0)
	}

	if iat, ok := claims["iat"].(float64); ok {
		payload.IssuedAt = time.Unix(int64(iat), 0)
	}

	return payload, nil
}
