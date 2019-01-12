package db

import (
	"database/sql"
	"fmt"
	"sesame/config"

	_ "github.com/lib/pq"
)

var Database *sql.DB

func init() {
	sslmode := "disable"
	Database, _ = sql.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=%s",
			config.DBUser, config.DBName, sslmode))
}