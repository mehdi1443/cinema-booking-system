package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

type Handlers struct {
	Auth handler.AuthHandler
	User handler.UserHandler
}

func SetUpRouter(handler *Handlers, config *bootstrap.Config, limiter *middleware.RedisLimiter) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS
	r.Use(middleware.SecureCORSMiddleware([]string{
		fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port),
	}))

	//HACK : this is just a exm for how to use limiter
	// r.Use(limiter.PerSecond(config.Constant.Limiter.HighRate))
	//FIXME : ADD AUTHORIZATION MIDDLEWARE!!!!!  exm :
	// r.Use(middleware.ControlUserAccess(/*// FIXME permision list*/))

	// v1 := r.Group("/v1")

	// v1.POST("/send-otp", handler.Auth.SendOTP)
	// v1.POST("/verify-otp", handler.Auth.VerifyOTP)
	// v1.POST("/register", handler.Auth.Register)
	// v1.POST("/login", handler.Auth.Login)
	// v1.POST("/refresh", handler.Auth.Refresh)

	// //FIXME : ADD AUTH MIDDLEWARE!!!!!
	// v1.POST("/logout", handler.Auth.Logout)
	return r
}
