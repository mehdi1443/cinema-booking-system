package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func JWTAuthentication(tokenService service_contract.JWTServiceInterface, appLogger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		payload, err := tokenService.ValidateToken(c.Request.Context(), tokenString, string(domainjwt.AccessTokenTag))
		if err != nil {
			if errors.Is(err, exception.ErrTokenExpired) {
				response.Error(c, http.StatusUnauthorized, "access token has expired")
				c.Abort()
				return
			}

			appLogger.Warn(logger.Auth, logger.JWT, "failed token validation", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, "token is malformed or invalid")
			c.Abort()
			return
		}

		if payload.TokenType != domainjwt.AccessTokenTag {
			response.Error(c, http.StatusForbidden, "invalid_token_type")
			c.Abort()
			return
		}

		isBlacklisted, err := tokenService.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil || (isBlacklisted != nil && *isBlacklisted) {
			appLogger.Warn(logger.Auth, logger.JWT, "blacklisted token execution attempt", map[logger.ExtraKey]interface{}{
				logger.UserID: payload.UserID,
			})
			response.Error(c, http.StatusUnauthorized, "token has been revoked")
			c.Abort()
			return
		}

		c.Set(domainjwt.CtxUserIDKey, payload.UserID)
		c.Set(domainjwt.CtxTheTokenKey, tokenString)
		c.Set(domainjwt.CtxTheTokenExpireAtKey, payload.ExpiresAt)
		c.Set(domainjwt.CtxPayloadKey, payload) // i don't know may be need it in future :)

		c.Next()
	}
}
