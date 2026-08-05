package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/application/service"
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

// TODO : REVIEW AND DESISIION MAKING ABOUT PROJECT STRUCTURE
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
	accessRSA, err := rsa.NewRSAKeys(config.Environment.JWT.AccessPrivateKey)
	if err != nil {
		log.Fatalf("build-access-rsa:WTF")
	}

	refreshRSA, err := rsa.NewRSAKeys(config.Environment.JWT.RefreshPrivateKey)
	if err != nil {
		log.Fatalf("build-refresh-rsa:WTF")
	}

	// --- HASHER ---
	bcryptHasher := hash.NewBcryptHasher()

	// --- REPOSITORIES ---
	userRepo := postgres_repo.NewUserRepositoryPostgres(pgDB)
	walletRepo := postgres_repo.NewWalletRepositoryPostgres(pgDB)
	rbacRepo := postgres_repo.NewRABCRepositoryPostgres(pgDB)
	otpRepo := redis_repo.NewOTPRepository(redisClient)
	refreshtokenRepo := redis_repo.NewRefreshTokenRepository(redisClient)
	// userpermissionRepo := redis_repo.NewUserPermissionCacheRepository(redisClient) //TODO USE IT!!!

	// --- SEEDER ---
	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, rbacRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Database seeded successfully!")

	// --- SERVICES ---
	jwtSrv := service.NewTokenService(accessRSA, refreshRSA, &config.Environment.JWT)
	smsSrv := service.NewSMSService()
	userSrv := service.NewUserService(userRepo)
	authSrv := service.NewAuthService(userRepo, otpRepo, refreshtokenRepo, jwtSrv, random.NewRandomGeneratorEngine(), &config.Environment.OTP, &config.Environment.JWT, smsSrv, hash.NewBcryptHasher())
	walletSrv := service.NewWalletService(userRepo, walletRepo)
	rbacSrv := service.NewRBACService(rbacRepo, userRepo)
	contextokenSrv := service.NewContextTokenService()

	services := router.Services{
		TokenService: jwtSrv,
	}

	// --- HANDLERS ---
	authHandler := *handler.NewAuthHandler(authSrv)
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
	//TODO : Im not sure is this clean or not(that import middleware in main and routes package both)
	limiter := middleware.NewRedisLimiter(redisClient)

	// --- ROUTES ---
	Engine := router.SetUpRouter(&handlers, &services, config, limiter)

	log.Printf("server running on %s:%s\n", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port), Engine))
}
