package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp" // ایمپورت برای اندپوینت متریک

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
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
	Comment   handler.CommentHandler
	Rating    handler.RatingHandler
	Search    handler.SearchHandler
	Booking   handler.BookingHandler
	Payment   handler.PaymentHandler
	Ticket    handler.TicketHandler
}

type Repositories struct {
	UserRepo                DB.UserInterface
	RBACRepo                DB.RBACInterface
	UserPermissionCacheRepo cache.UserPermissionCacheInterface
	CinemaRepo              DB.CinemaInterface
	BookingRepo             DB.BookingRepositoryInterface
}

type Services struct {
	TokenService        service_contract.JWTServiceInterface
	RBACService         service_contract.RBACServiceInterface
	ContextTokenService service_contract.ContextTokenServiceInterface
	UserService         service_contract.UserServiceInterface
	CaptchaService      service_contract.CaptchaServiceInterface
}

func SetUpRouter(
	handler *Handlers,
	services *Services,
	repositories *Repositories,
	config *bootstrap.Config,
	limiter *middleware.RedisLimiter,
	appLogger domainlogger.Logger,
	promMiddleware *middleware.PrometheusMiddleware,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.GinLoggerMiddleware(appLogger))
	r.Use(gin.Recovery())

	r.Use(promMiddleware.Handle())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(middleware.SecureCORSMiddleware([]string{
		fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port),
	}))

	captchaMiddleware := middleware.NewCaptchaSolving(services.CaptchaService, &config.Environment.DebugModeOptions, appLogger)

	v1 := r.Group("/v1")
	v1.Use(limiter.PerMinute(config.Environment.Limiter.HighRate))
	v1.Use(middleware.RouteAuthorization(
		services.TokenService,
		repositories.RBACRepo,
		repositories.UserPermissionCacheRepo,
		services.ContextTokenService,
		repositories.UserRepo,
		repositories.CinemaRepo,
		repositories.BookingRepo,
		&config.Environment.UserPermissionCache,
		appLogger,
	))

	SetUpUserRouter(
		v1,
		&handler.User,
		&handler.Wallet,
	)

	SetUpAuthRouter(
		v1,
		&handler.Auth,
		captchaMiddleware,
	)

	SetUpRBACRouter(
		v1,
		&handler.RBAC,
	)

	SetUpGenreRouter(v1, &handler.Genre)

	SetUpPersonRouter(v1, &handler.Person)

	SetUpMovieRouter(v1, &handler.Movie)

	SetUpAddressRouter(v1, &handler.Address)

	SetUpFacilityRouter(v1, &handler.Facilitie)

	SetUpCinemaRouter(v1, &handler.Cinema)

	SetUpFavoriteRouter(v1, &handler.Favorite)

	SetUpCommentRouter(v1, &handler.Comment)

	SetUpRatingRouter(v1, &handler.Rating)

	SetUpCaptchaRouter(v1, &handler.Captcha)

	SetUpSearchRouter(v1, &handler.Search)

	SetUpBookingRouter(v1, &handler.Booking)

	SetUpPaymentRouter(v1, &handler.Payment)

	SetUpTicketRouter(v1, &handler.Ticket)
	if err := middleware.ValidateRoutePolicies(r.Routes()); err != nil {
		panic(err)
	}

	return r
}
