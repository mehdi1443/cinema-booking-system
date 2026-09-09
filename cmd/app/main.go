package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/application/service"
	"github.com/mehdi1443/cinema-booking-system/internal/application/validation"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/broker"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	infralogger "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/logger"
	infralmetrics "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/metrics"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/MinIO"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/elasticsearch"
	postgres_repo "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/postgres"
	redis_repo "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/redis"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/seed"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/router"
	"github.com/mehdi1443/cinema-booking-system/pkg/captcha"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
	"github.com/mehdi1443/cinema-booking-system/pkg/random"
	"github.com/mehdi1443/cinema-booking-system/pkg/rsa"
	mail "github.com/mehdi1443/cinema-booking-system/pkg/wneessen-go-mail"
	"github.com/mehdi1443/cinema-booking-system/static"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	config := bootstrap.Run()

	appLogger := infralogger.New(infralogger.Config{
		ServiceName: "cinema-booking-system",
		Environment: config.Environment.Logger.Level,
		Level:       config.Environment.Logger.Level,
		AddSource:   true,
	})

	appLogger.Info(ctx, "starting application bootstrapping")

	pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database, appLogger)
	redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)
	storage, err := MinIO.NewS3Storage(config.Environment)
	if err != nil {
		appLogger.Error(ctx, "failed to initialize S3 storage", domainlogger.Error(err))
		os.Exit(1)
	}

	esClient, err := elasticsearch.NewClient(config.Environment.Elasticsearch.URL)
	if err != nil {
		appLogger.Error(ctx, "failed to initialize Elasticsearch", domainlogger.Error(err))
		os.Exit(1)
	}
	searchRepo := elasticsearch.NewSearchRepository(esClient)

	rabbitConn, err := amqp.Dial(config.Environment.RabbitMQ.URL)
	if err != nil {
		appLogger.Error(ctx, "failed to connect to RabbitMQ", domainlogger.Error(err))
		os.Exit(1)
	}
	defer rabbitConn.Close()

	rabbitPubCh, err := rabbitConn.Channel()
	if err != nil {
		appLogger.Error(ctx, "failed to open RabbitMQ publisher channel", domainlogger.Error(err))
		os.Exit(1)
	}
	defer rabbitPubCh.Close()

	searchPub, err := broker.NewSearchEventPublisher(rabbitPubCh)
	if err != nil {
		appLogger.Error(ctx, "failed to create search event publisher", domainlogger.Error(err))
		os.Exit(1)
	}

	rabbitWorkerCh, err := rabbitConn.Channel()
	if err != nil {
		appLogger.Error(ctx, "failed to open RabbitMQ worker channel", domainlogger.Error(err))
		os.Exit(1)
	}
	defer rabbitWorkerCh.Close()

	searchWorker, err := broker.NewSearchWorker(rabbitWorkerCh, searchRepo, appLogger)
	if err != nil {
		appLogger.Error(ctx, "failed to create search worker", domainlogger.Error(err))
		os.Exit(1)
	}

	go func() {
		if err := searchWorker.Start(ctx); err != nil {
			appLogger.Error(ctx, "Search worker stopped with error", domainlogger.Error(err))
		}
	}()

	if pgDB != nil && redisClient != nil {
		appLogger.Info(ctx, "Postgres, Redis, and Object Storage initialized successfully")
	}

	accessKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateAccessPath, config.Environment.PEM.PublicAccessPath, 2048)
	if err != nil {
		appLogger.Error(ctx, "failed to initialize access RSA keys", domainlogger.Error(err))
		os.Exit(1)
	}

	refreshKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateRefreshPath, config.Environment.PEM.PublicRefreshPath, 2048)
	if err != nil {
		appLogger.Error(ctx, "failed to initialize refresh RSA keys", domainlogger.Error(err))
		os.Exit(1)
	}

	bcryptHasher := hash.NewBcryptHasher()

	statics, err := static.InitStaticFiles()
	if err != nil {
		appLogger.Error(ctx, "failed to initialize static files", domainlogger.Error(err))
		os.Exit(1)
	}

	captchaMan, err := captcha.NewCaptchaMan(config.Environment.Captcha.Type)
	if err != nil {
		appLogger.Error(ctx, "failed to initialize Captcha Man", domainlogger.Error(err))
		os.Exit(1)
	}

	userRepo := postgres_repo.NewUserRepositoryPostgres(pgDB)
	walletRepo := postgres_repo.NewWalletRepositoryPostgres(pgDB)
	rbacRepo := postgres_repo.NewRABCRepositoryPostgres(pgDB)
	otpRepo := redis_repo.NewOTPRepository(redisClient)
	refreshtokenRepo := redis_repo.NewRefreshTokenRepository(redisClient)
	userpermissionRepo := redis_repo.NewUserPermissionCacheRepository(redisClient)
	blacklistRepo := redis_repo.NewTokenBlacklistRedis(redisClient)
	genreRepo := postgres_repo.NewGenreRepository(pgDB)
	personRepo := postgres_repo.NewPersonRepository(pgDB)
	movieRepo := postgres_repo.NewMovieRepository(pgDB)
	captchaRepo := redis_repo.NewCaptchaRepositoryRedis(redisClient)
	provinceRepo := postgres_repo.NewProvinceRepository(pgDB)
	cityRepo := postgres_repo.NewCityRepository(pgDB)
	facilityRepo := postgres_repo.NewFacilityRepository(pgDB)
	cinemaRepo := postgres_repo.NewCinemaRepository(pgDB)
	favoriteRepo := postgres_repo.NewFavoriteRepository(pgDB)
	commentRepo := postgres_repo.NewCommentRepository(pgDB)
	ratingRepo := postgres_repo.NewRatingRepository(pgDB)
	bookingRepo := postgres_repo.NewBookingRepository(pgDB)
	paymentRepo := postgres_repo.NewPaymentRepository(pgDB)
	ticketRepo := postgres_repo.NewTicketRepository(pgDB)

	repositories := router.Repositories{
		UserRepo:                userRepo,
		RBACRepo:                rbacRepo,
		UserPermissionCacheRepo: userpermissionRepo,
		CinemaRepo:              cinemaRepo,
		BookingRepo:             bookingRepo,
	}

	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, rbacRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		// RBAC and the bootstrap administrator are required for a usable API.
		// Starting without them would leave every protected endpoint unusable.
		appLogger.Error(ctx, "failed to initialize roles and bootstrap administrator", domainlogger.Error(err))
		os.Exit(1)
	} else {
		appLogger.Info(ctx, "database roles seeded successfully")
	}

	addressSeeder := seed.NewAddressSeeder(provinceRepo, cityRepo)
	if err := addressSeeder.Seed(ctx); err != nil {
		appLogger.Warn(ctx, "address seeder finished with potential existing entries or non-fatal errors", domainlogger.Error(err))
	} else {
		appLogger.Info(ctx, "database provinces and cities seeded successfully")
	}

	emailSender := mail.NewMailer(mail.EmailConfig{
		From:     config.Environment.Email.CinemaTicketEmail,
		Password: config.Environment.Email.CinemaTicketEmailAppPassword,
		SMTPHost: config.Environment.Email.SMTPHost,
		SMTPPort: config.Environment.Email.SMTPPort,
	})

	jwtSrv := service.NewTokenService(accessKeys, refreshKeys, &config.Environment.JWT, blacklistRepo, appLogger)
	smsSrv := service.NewSMSService()
	emailSrv := service.NewEmailService(emailSender, &config.Environment.DebugModeOptions, statics, appLogger)
	userSrv := service.NewUserService(userRepo, bcryptHasher, appLogger)
	authSrv := service.NewAuthServiceWithDebug(userRepo, otpRepo, refreshtokenRepo, jwtSrv, random.NewRandomGeneratorEngine(), &config.Environment.OTP, &config.Environment.JWT, &config.Environment.DebugModeOptions, smsSrv, emailSrv, hash.NewBcryptHasher(), appLogger)
	walletSrv := service.NewWalletService(userRepo, walletRepo, appLogger)
	rbacSrv := service.NewRBACService(rbacRepo, userRepo, userpermissionRepo, &config.Environment.UserPermissionCache, appLogger)
	contextokenSrv := service.NewContextTokenService()

	genreSrv := service.NewGenreService(genreRepo, appLogger)
	personSrv := service.NewPersonService(personRepo, storage, appLogger, searchPub)
	movieSrv := service.NewMovieService(movieRepo, storage, appLogger, searchPub)
	captchaSrv := service.NewCaptchaService(captchaRepo, &config.Environment.Captcha, captchaMan, appLogger)
	addressSrv := service.NewAddressService(provinceRepo, cityRepo, appLogger)
	facilitySrv := service.NewFacilityService(facilityRepo, appLogger)
	cinemaSrv := service.NewCinemaService(cinemaRepo, storage, appLogger, searchPub)
	favoriteSrv := service.NewFavoriteService(favoriteRepo, appLogger)
	commentSrv := service.NewCommentService(commentRepo, appLogger)
	ratingSrv := service.NewRatingService(ratingRepo, movieRepo, cinemaRepo, appLogger)
	searchSrv := service.NewSearchService(searchRepo, movieSrv, cinemaSrv, personSrv, appLogger)
	hallSrv := service.NewHallService(bookingRepo, cinemaRepo, appLogger)
	showtimeSrv := service.NewShowtimeService(bookingRepo, bookingRepo, movieRepo, bookingRepo, appLogger)
	reservationSrv := service.NewReservationService(bookingRepo, &config.Environment.Reservation, appLogger)
	ticketSrv := service.NewTicketService(ticketRepo, storage, appLogger)
	paymentSrv := service.NewPaymentService(paymentRepo, bookingRepo, broker.NewMockPaymentGateway(), ticketSrv, config.Environment.Reservation.RefundCutoff(), appLogger)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				if err := reservationSrv.ExpirePending(ctx, now); err != nil {
					appLogger.Error(ctx, "failed to expire pending reservations", domainlogger.Error(err))
				}
			}
		}
	}()

	services := router.Services{
		TokenService:        jwtSrv,
		RBACService:         rbacSrv,
		ContextTokenService: contextokenSrv,
		UserService:         userSrv,
		CaptchaService:      captchaSrv,
	}

	authHandler := *handler.NewAuthHandler(authSrv, contextokenSrv, appLogger)
	userHandler := *handler.NewUserHandler(userSrv, rbacSrv, contextokenSrv, appLogger)
	rbacHandler := *handler.NewRBACHandler(rbacSrv, contextokenSrv, appLogger)
	walletHandler := *handler.NewWalletHandler(walletSrv, contextokenSrv, appLogger)

	genreHandler := *handler.NewGenreHandler(genreSrv, appLogger)
	personHandler := *handler.NewPersonHandler(personSrv, appLogger)
	movieHandler := *handler.NewMovieHandler(movieSrv, appLogger)
	captchaHandler := *handler.NewCaptchaHandler(captchaSrv, contextokenSrv, appLogger)
	addressHandler := *handler.NewAddressHandler(addressSrv, appLogger)
	facilitiesHandler := *handler.NewFacilityHandler(facilitySrv, appLogger)
	cinemaHandler := *handler.NewCinemaHandler(cinemaSrv, appLogger)
	favoriteHandler := *handler.NewFavoriteHandler(favoriteSrv, contextokenSrv, appLogger)
	commentHandler := *handler.NewCommentHandler(commentSrv, appLogger)
	ratingHandler := *handler.NewRatingHandler(ratingSrv, appLogger)
	searchHandler := *handler.NewSearchHandler(searchSrv, appLogger)
	bookingHandler := *handler.NewBookingHandler(hallSrv, showtimeSrv, reservationSrv, appLogger)
	paymentHandler := *handler.NewPaymentHandler(paymentSrv, appLogger)
	ticketHandler := *handler.NewTicketHandler(ticketSrv, appLogger)

	handlers := router.Handlers{
		Auth:      authHandler,
		User:      userHandler,
		RBAC:      rbacHandler,
		Wallet:    walletHandler,
		Genre:     genreHandler,
		Person:    personHandler,
		Movie:     movieHandler,
		Captcha:   captchaHandler,
		Address:   addressHandler,
		Facilitie: facilitiesHandler,
		Cinema:    cinemaHandler,
		Favorite:  favoriteHandler,
		Comment:   commentHandler,
		Rating:    ratingHandler,
		Search:    searchHandler,
		Booking:   bookingHandler,
		Payment:   paymentHandler,
		Ticket:    ticketHandler,
	}

	limiter := middleware.NewRedisLimiter(redisClient, appLogger)

	promClient := infralmetrics.NewPrometheusMetrics()
	promMiddleware := middleware.NewPrometheusMiddleware(promClient)

	Engine := router.SetUpRouter(&handlers, &services, &repositories, config, limiter, appLogger, promMiddleware)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("iranian-phone", validation.ValidateIranianPhone)
		v.RegisterValidation("otp-code", validation.ValidateOTPCode(int(config.Environment.OTP.Digit)))
		v.RegisterValidation("user-password", validation.ValidatePassword(&config.Environment.DebugModeOptions))
		v.RegisterValidation("role-title", validation.ValidateRoleTitle)
		v.RegisterValidation("notblank", validation.ValidateNotBlank)
	}

	serverAddr := fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port)
	appLogger.Info(ctx, fmt.Sprintf("HTTP server starting on %s", serverAddr))

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           Engine,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error(ctx, "HTTP server failed to listen and serve", domainlogger.Error(err))
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			appLogger.Error(context.Background(), "HTTP server graceful shutdown failed", domainlogger.Error(err))
			os.Exit(1)
		}
	}
}
