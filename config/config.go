package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	DBName string
	DBUser string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DBName = os.Getenv("DATABASE_NAME")
	DBUser = os.Getenv("DATABASE_USER")
}
