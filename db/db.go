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

	err := Database.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, _ := Database.Query("SELECT name FROM cities WHERE id = 1")

	var n string
	rows.Scan(&n)
	fmt.Println(n)
}
