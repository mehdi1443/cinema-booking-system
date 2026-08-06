package main

import (
	"context"
	"fmt"
	"log"
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
	// --- LOAD CONFIG & CONST
	ctx := context.Background()
	config := bootstrap.Run()
	log.Println(ctx)
	log.Println(config)

	// --- DATABASES CONNECTIONS ---
	pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)
	storage, err := MinIO.NewS3Storage(config.Environment)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	if pgDB != nil && redisClient != nil && storage != nil {
		fmt.Println("Database and Redis connected successfully!")
	}

	// --- RSA KEYS ---
	accessKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateAccessPath, config.Environment.PEM.PublicAccessPath, 2048)
	if err != nil {
		log.Fatalf("build-access-rsa:WTF:Key initialization failed: %v", err)
	}
	fmt.Println("build-access-rsa:ok")

	refreshKeys, err := rsa.LoadOrCreate(config.Environment.PEM.PrivateRefreshPath, config.Environment.PEM.PublicRefreshPath, 2048)
	if err != nil {
		log.Fatalf("build-refresh-rsa:WTF:Key initialization failed: %v", err)
	}
	fmt.Println("build-refresh-rsa:ok")

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

	repositories := router.Repositories{
		UserRepo:                userRepo,
		RBACRepo:                rbacRepo,
		UserPermissionCacheRepo: userpermissionRepo,
	}

	// --- SEEDER ---
	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, rbacRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		log.Printf("Failed to run seeders: %v\n", err) // TODO : NEED EDIT WHEN USER ALREADY EXISTS
	}

	log.Println("Database seeded successfully!")

	// --- SERVICES ---
	jwtSrv := service.NewTokenService(accessKeys, refreshKeys, &config.Environment.JWT, blacklistRepo)
	smsSrv := service.NewSMSService()
	emailSrv := service.NewEmailService()
	userSrv := service.NewUserService(userRepo, bcryptHasher)
	authSrv := service.NewAuthService(userRepo, otpRepo, refreshtokenRepo, jwtSrv, random.NewRandomGeneratorEngine(), &config.Environment.OTP, &config.Environment.JWT, smsSrv, emailSrv, hash.NewBcryptHasher())
	walletSrv := service.NewWalletService(userRepo, walletRepo)
	rbacSrv := service.NewRBACService(rbacRepo, userRepo, userpermissionRepo, &config.Environment.UserPermissionCache)
	contextokenSrv := service.NewContextTokenService()

	services := router.Services{
		TokenService:        jwtSrv,
		RBACService:         rbacSrv,
		ContextTokenService: contextokenSrv,
		UserService:         userSrv,
	}

	// --- HANDLERS ---
	authHandler := *handler.NewAuthHandler(authSrv, contextokenSrv)
	userHandler := *handler.NewUserHandler(userSrv, rbacSrv, contextokenSrv)
	rbacHandler := *handler.NewRBACHandler(rbacSrv, contextokenSrv)
	walletHandler := *handler.NewWalletHandler(walletSrv, contextokenSrv)
	handlers := router.Handlers{
		Auth:   authHandler,
		User:   userHandler,
		RBAC:   rbacHandler,
		Wallet: walletHandler,
	}

	// --- LIMITER ---
	limiter := middleware.NewRedisLimiter(redisClient)

	// --- ROUTES ---
	Engine := router.SetUpRouter(&handlers, &services, &repositories, config, limiter)

	// --- VALIDATIONS ---
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("iranian-phone", validation.ValidateIranianPhone)
		v.RegisterValidation("otp-code", validation.ValidateOTPCode(int(config.Environment.OTP.Digit)))
		v.RegisterValidation("user-password", validation.ValidatePassword(&config.Environment.DebugModeOptions))
		//HACK : gin handle email validation by default
		// 	v.RegisterValidation("email", bookableDate)
		v.RegisterValidation("role-title", validation.ValidateRoleTitle)
	}

	log.Printf("server running on %s:%s\n", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port), Engine))
}
