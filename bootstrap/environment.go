package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	PDB DataBase
	RDB Redis
}

type DataBase struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Redis struct {
	Host      string
	Port      string
	Password  string
	RDBNumber string
}

func NewEnvironment() *Environment {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	return &Environment{
		PDB: DataBase{
			Host:     os.Getenv("PDB_HOST"),
			Port:     os.Getenv("PDB_PORT"),
			User:     os.Getenv("PDB_USER"),
			Password: os.Getenv("PDB_PASSWORD"),
			Name:     os.Getenv("PDB_NAME"),
		},
		RDB: Redis{
			Host:      os.Getenv("RDB_HOST"),
			Port:      os.Getenv("RDB_PORT"),
			Password:  os.Getenv("RDB_PASSWORD"),
			RDBNumber: os.Getenv("RDB_NUMBER"),
		},
	}
}
