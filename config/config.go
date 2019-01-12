package config

import "os"

var (
	DBName = os.Getenv("DATABASE_NAME")
	DBUser = os.Getenv("DATABASE_USER")
)
