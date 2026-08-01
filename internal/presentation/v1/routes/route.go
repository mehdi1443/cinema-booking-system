package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/controller"
)

type Controllers struct {
	Auth controller.AuthController
	User controller.UserController
}

func SetUpRouter(handler *Controllers, config *bootstrap.Config, limiter *middleware.RedisLimiter) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS
	r.Use(middleware.SecureCORSMiddleware([]string{
		fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port),
	}))

	//HACK : this is just a exm for how to use limiter
	// r.Use(limiter.PerSecond(config.Constant.Limiter.HighRate))

	v1 := r.Group("/v1")

	v1.POST("/send-otp", handler.Auth.SendOTP)
	v1.POST("/verify-otp", handler.Auth.VerifyOTP)
	v1.POST("/register", handler.Auth.Register)
	v1.POST("/login", handler.Auth.Login)
	v1.POST("/refresh", handler.Auth.Refresh)

	//FIXME : ADD AUTH MIDDLEWARE!!!!!
	v1.POST("/logout", handler.Auth.Logout)
	return r
}
