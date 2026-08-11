package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

type Handlers struct {
	Auth      handler.AuthHandler
	User      handler.UserHandler
	RBAC      handler.RBACHandler
	Wallet    handler.WalletHandler
	Genre     handler.GenreHandler
	Person    handler.PersonHandler
	Movie     handler.MovieHandler
	Captcha   handler.CaptchaHandler
	Address   handler.AddressHandler
	Facilitie handler.FacilityHandler
	Cinema    handler.CinemaHandler
	Favorite  handler.FavoriteHandler
}

type Repositories struct {
	UserRepo                DB.UserInterface
	RBACRepo                DB.RBACInterface
	UserPermissionCacheRepo cache.UserPermissionCacheInterface
}

type Services struct {
	TokenService        service_contract.JWTServiceInterface
	RBACService         service_contract.RBACServiceInterface
	ContextTokenService service_contract.ContextTokenServiceInterface
	UserService         service_contract.UserServiceInterface
	CaptchaService      service_contract.CaptchaServiceInterface
}

func SetUpRouter(handler *Handlers, services *Services, repositories *Repositories, config *bootstrap.Config, limiter *middleware.RedisLimiter, appLogger logger.Logger) *gin.Engine {
	r := gin.New()

	// gin logger
	r.Use(gin.Logger())
	// gin panic recovery
	r.Use(gin.Recovery())

	// CORS
	r.Use(middleware.SecureCORSMiddleware([]string{
		fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port),
	}))

	authMiddleware := middleware.JWTAuthentication(services.TokenService, appLogger)
	captchaMiddleware := middleware.NewCaptchaSolving(services.CaptchaService, &config.Environment.DebugModeOptions, appLogger)

	v1 := r.Group("/v1")

	SetUpUserRouter(
		v1,
		&handler.User,
		&handler.RBAC,
		&handler.Wallet,
		authMiddleware,
		limiter,
		&config.Environment.Limiter,
		repositories.RBACRepo,
		repositories.UserPermissionCacheRepo,
		services.ContextTokenService,
		repositories.UserRepo,
		&config.Environment.UserPermissionCache,
		appLogger,
	)

	SetUpAuthRouter(
		v1,
		&handler.Auth,
		authMiddleware,
		limiter,
		&config.Environment.Limiter,
		repositories.RBACRepo,
		repositories.UserPermissionCacheRepo,
		services.ContextTokenService,
		repositories.UserRepo,
		&config.Environment.UserPermissionCache,
		appLogger,
		captchaMiddleware,
	)

	SetUpRBACRouter(
		v1,
		&handler.RBAC,
		authMiddleware,
		limiter,
		&config.Environment.Limiter,
		repositories.RBACRepo,
		repositories.UserPermissionCacheRepo,
		services.ContextTokenService,
		repositories.UserRepo,
		&config.Environment.UserPermissionCache,
		appLogger,
	)

	SetUpGenreRouter(v1, &handler.Genre)

	SetUpPersonRouter(v1, &handler.Person)

	SetUpMovieRouter(v1, &handler.Movie)

	SetUpAddressRouter(v1, &handler.Address)

	SetUpFacilityRouter(v1, &handler.Facilitie)

	SetUpCinemaRouter(v1, &handler.Cinema)

	SetUpFavoriteRouter(
		v1,
		&handler.Favorite,
		authMiddleware,
		repositories.RBACRepo,
		repositories.UserPermissionCacheRepo,
		services.ContextTokenService,
		repositories.UserRepo,
		&config.Environment.UserPermissionCache,
		appLogger,
	)

	SetUpCaptchaRouter(v1, &handler.Captcha, limiter, &config.Environment.Limiter, appLogger)
	return r
}
