package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/postgres"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/seed"
	"github.com/mehdi1443/cinema-booking-system/pkg/hasher"
)

func main() {
	ctx := context.Background()
	config := bootstrap.Run()

	pgDB := database.NewPostgresDatabase(&config.Environment.PDB, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.RDB)

	if pgDB != nil && redisClient != nil {
		fmt.Println("Database and Redis connected successfully!")
	}

	userRepo := postgres.NewUserRepositoryPostgres()
	bcryptHasher := hasher.NewBcryptHasher()

	seeder := seed.NewRoleSeeder(config.Constant.User, userRepo, pgDB, bcryptHasher)
	if err := seeder.Seed(ctx); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Database seeded successfully!")
}
