package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func JWTAuthentication(tokenService service_contract.JWTServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization_header_required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid_authorization_format")
			return
		}

		tokenString := parts[1]

		payload, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, exception.ErrTokenExpired) {
				response.Error(c, http.StatusUnauthorized, "access token has expired")
				return
			}

			response.Error(c, http.StatusUnauthorized, "token is malformed or invalid")
			return
		}

		if payload.TokenType != domainjwt.AccessTokenTag {
			response.Error(c, http.StatusForbidden, "invalid_token_type")
			return
		}

		c.Set(domainjwt.CtxUserIDKey, payload.UserID)
		c.Set(domainjwt.CtxPayloadKey, payload) // i don't know may be need it in future :)

		c.Next()
	}
}
