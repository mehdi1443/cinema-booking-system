package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/application/service"
	"github.com/mehdi1443/cinema-booking-system/internal/application/validation"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/MinIO"
	postgres_repo "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/postgres"
	redis_repo "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/redis"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/seed"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/router"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/random"
	"github.com/mehdi1443/cinema-booking-system/pkg/rsa"
)

// comment types
// TODO
// BUG
// FIXME
// HACK
// XXX
// [ ]
// [x]
// !
// ?

func main() {
	// --- LOAD CONFIG & CONST ---
	ctx := context.Background()
	config := bootstrap.Run()

	// --- LOGGER INITIALIZATION ---
	appLogger := logger.NewLogger(&config.Environment.Logger)
	appLogger.Info(logger.General, logger.Startup, "starting application bootstrapping", nil)

	// --- DATABASES CONNECTIONS ---
	pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)
	storage, err := MinIO.NewS3Storage(config.Environment)
	if err != nil {
		appLogger.Fatal(logger.General, logger.Startup, "failed to initialize S3 storage", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}

	if pgDB != nil && redisClient != nil && storage != nil {
		appLogger.Info(logger.General, logger.Startup, "Postgres, Redis, and Object Storage initialized successfully", nil)
	}

	// --- RSA KEYS ---
	accessKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateAccessPath, config.Environment.PEM.PublicAccessPath, 2048)
	if err != nil {
		appLogger.Fatal(logger.General, logger.Startup, "failed to initialize access RSA keys", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}

	refreshKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateRefreshPath, config.Environment.PEM.PublicRefreshPath, 2048)
	if err != nil {
		appLogger.Fatal(logger.General, logger.Startup, "failed to initialize refresh RSA keys", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}

	// --- HASHER ---
	bcryptHasher := hash.NewBcryptHasher()

	// --- REPOSITORIES ---
	userRepo := postgres_repo.NewUserRepositoryPostgres(pgDB)
	walletRepo := postgres_repo.NewWalletRepositoryPostgres(pgDB)
	rbacRepo := postgres_repo.NewRABCRepositoryPostgres(pgDB)
	otpRepo := redis_repo.NewOTPRepository(redisClient)
	refreshtokenRepo := redis_repo.NewRefreshTokenRepository(redisClient)
	userpermissionRepo := redis_repo.NewUserPermissionCacheRepository(redisClient) //TODO USE IT!! Actually Test it
	blacklistRepo := redis_repo.NewTokenBlacklistRedis(redisClient)
	genreRepo := postgres_repo.NewGenreRepository(pgDB)
	personRepo := postgres_repo.NewPersonRepository(pgDB)
	movieRepo := postgres_repo.NewMovieRepository(pgDB)

	repositories := router.Repositories{
		UserRepo:                userRepo,
		RBACRepo:                rbacRepo,
		UserPermissionCacheRepo: userpermissionRepo,
	}

	// --- SEEDER ---
	// TODO : NEED EDIT WHEN USER ALREADY EXISTS
	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, rbacRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		appLogger.Warn(logger.General, logger.Startup, "role seeder finished with potential existing entries or non-fatal errors", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	} else {
		appLogger.Info(logger.General, logger.Startup, "database roles seeded successfully", nil)
	}

	// --- SERVICES ---
	jwtSrv := service.NewTokenService(accessKeys, refreshKeys, &config.Environment.JWT, blacklistRepo, appLogger)
	smsSrv := service.NewSMSService()
	emailSrv := service.NewEmailService()
	userSrv := service.NewUserService(userRepo, bcryptHasher, appLogger)
	authSrv := service.NewAuthService(userRepo, otpRepo, refreshtokenRepo, jwtSrv, random.NewRandomGeneratorEngine(), &config.Environment.OTP, &config.Environment.JWT, smsSrv, emailSrv, hash.NewBcryptHasher(), appLogger)
	walletSrv := service.NewWalletService(userRepo, walletRepo, appLogger)
	rbacSrv := service.NewRBACService(rbacRepo, userRepo, userpermissionRepo, &config.Environment.UserPermissionCache, appLogger)
	contextokenSrv := service.NewContextTokenService()
	genreSrv := service.NewGenreService(genreRepo)
	personSrv := service.NewPersonService(personRepo, storage)
	movieSrv := service.NewMovieService(movieRepo, storage)

	services := router.Services{
		TokenService:        jwtSrv,
		RBACService:         rbacSrv,
		ContextTokenService: contextokenSrv,
		UserService:         userSrv,
	}

	// --- HANDLERS ---
	authHandler := *handler.NewAuthHandler(authSrv, contextokenSrv, appLogger)
	userHandler := *handler.NewUserHandler(userSrv, rbacSrv, contextokenSrv, appLogger)
	rbacHandler := *handler.NewRBACHandler(rbacSrv, contextokenSrv, appLogger)
	walletHandler := *handler.NewWalletHandler(walletSrv, contextokenSrv, appLogger)
	genreHandler := *handler.NewGenreHandler(genreSrv)
	personHandler := *handler.NewPersonHandler(personSrv)
	movieHandler := *handler.NewMovieHandler(movieSrv)

	handlers := router.Handlers{
		Auth:   authHandler,
		User:   userHandler,
		RBAC:   rbacHandler,
		Wallet: walletHandler,
		Genre:  genreHandler,
		Person: personHandler,
		Movie:  movieHandler,
	}

	// --- LIMITER ---
	limiter := middleware.NewRedisLimiter(redisClient, appLogger)

	// --- ROUTES ---
	Engine := router.SetUpRouter(&handlers, &services, &repositories, config, limiter, appLogger)

	// --- VALIDATIONS ---
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("iranian-phone", validation.ValidateIranianPhone)
		v.RegisterValidation("otp-code", validation.ValidateOTPCode(int(config.Environment.OTP.Digit)))
		v.RegisterValidation("user-password", validation.ValidatePassword(&config.Environment.DebugModeOptions))
		//HACK : gin handle email validation by default
		// 	v.RegisterValidation("email", bookableDate)
		v.RegisterValidation("role-title", validation.ValidateRoleTitle)
	}

	serverAddr := fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port)
	appLogger.Info(logger.General, logger.Startup, fmt.Sprintf("HTTP server starting on %s", serverAddr), nil)

	if err := http.ListenAndServe(serverAddr, Engine); err != nil {
		appLogger.Fatal(logger.General, logger.Startup, "HTTP server failed to listen and serve", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
	}
}
