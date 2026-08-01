package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

// FIXME : Response Format
func JWTAuthentication(tokenService service_contract.JWTServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization_header_required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_authorization_format",
			})
			return
		}

		tokenString := parts[1]

		payload, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, service_contract.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "token_expired",
					"message": "access token has expired, please use refresh token",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "token is malformed or invalid",
			})
			return
		}

		if payload.TokenType != domainjwt.AccessTokenTag {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid_token_type",
			})
			return
		}

		c.Set(domainjwt.CtxUserIDKey, payload.UserID)
		c.Set(domainjwt.CtxPayloadKey, payload) // i don't know may be need it in future :)

		c.Next()
	}
}
