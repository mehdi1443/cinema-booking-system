package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func JWTAuthentication(tokenService service_contract.JWTServiceInterface, appLogger domainlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization_header_required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid_authorization_format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		payload, err := tokenService.ValidateToken(ctx, tokenString, string(domainjwt.AccessTokenTag))
		if err != nil {
			if errors.Is(err, exception.ErrTokenExpired) {
				response.Error(c, http.StatusUnauthorized, "access token has expired")
				c.Abort()
				return
			}

			appLogger.Warn(ctx, "failed token validation", domainlogger.Error(err))
			response.Error(c, http.StatusUnauthorized, "token is malformed or invalid")
			c.Abort()
			return
		}

		if payload.TokenType != domainjwt.AccessTokenTag {
			response.Error(c, http.StatusForbidden, "invalid_token_type")
			c.Abort()
			return
		}

		isBlacklisted, err := tokenService.IsBlacklisted(ctx, tokenString)
		if err != nil || (isBlacklisted != nil && *isBlacklisted) {

			appLogger.Warn(ctx, "blacklisted token execution attempt", domainlogger.Int64("user_id", payload.UserID))
			response.Error(c, http.StatusUnauthorized, "token has been revoked")
			c.Abort()
			return
		}

		c.Set(domainjwt.CtxUserIDKey, payload.UserID)
		c.Set(domainjwt.CtxTheTokenKey, tokenString)
		c.Set(domainjwt.CtxTheTokenExpireAtKey, payload.ExpiresAt)
		c.Set(domainjwt.CtxPayloadKey, payload)

		c.Next()
	}
}
