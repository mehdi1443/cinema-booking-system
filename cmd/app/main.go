package main

import (
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
)

func main() {
	config := bootstrap.Run()

	pgDB := database.NewPostgresDatabase(&config.Environment.PDB, &config.Constant.Database)
	redisClient := database.NewRedisDatabase(&config.Environment.RDB)

	if pgDB != nil && redisClient != nil {
		fmt.Println("Database and Redis connected successfully!")
	}
}
