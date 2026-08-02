package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	postgres "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/postgres/user"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/seed"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/controller"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/routes"
	"github.com/mehdi1443/cinema-booking-system/pkg/hash"
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
	ctx := context.Background()
	config := bootstrap.Run()

	pgDB := database.NewPostgresDatabase(&config.Environment.DataBase, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.Redis, &config.Constant.Redis)

	//TODO : Im not sure is this clean or not(that import middleware in main and routes package both)
	limiter := middleware.NewRedisLimiter(redisClient)

	if pgDB != nil && redisClient != nil {
		fmt.Println("Database and Redis connected successfully!")
	}

	userRepo := postgres.NewUserRepositoryPostgres(pgDB)
	bcryptHasher := hash.NewBcryptHasher()

	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Database seeded successfully!")

	/******************************************************************************************************/

	// --- SERVICES ---

	// --- CONTROLLER ---
	auth_controller := controller.NewAuthController()
	user_controller := controller.NewUserController()

	// --- ROUTE ---
	controllers := &routes.Controllers{
		Auth: auth_controller,
		User: user_controller,
	}
	r := routes.SetUpRouter(controllers, config, limiter)

	log.Printf("server running on %s:%s\n", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Environment.HttpServer.Domain, config.Environment.HttpServer.Port), r))
}
