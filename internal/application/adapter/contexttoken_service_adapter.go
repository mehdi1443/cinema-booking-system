package service_contract

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type ContextTokenServiceInterface interface {
	GetTokenFromGinContext(gctx *gin.Context) (dto.TokenInContext, error)
	GetJWTFromGinContext(gctx *gin.Context) (dto.JWTTokenInContext, error)
}
