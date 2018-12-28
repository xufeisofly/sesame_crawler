package main

import (
	"database/sql"
	"fmt"
	"sesame/dao"
	"sesame/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
)

func main() {
	app := iris.Default()

	db, _ := sql.Open("mysql", "norris@/sesame")
	cityDao := dao.CityDAO{db}
	defer cityDao.Close()

	city := cityDao.Get(1)
	fmt.Println(city.Name)

	app.Get("/tickets", handler.TicketList)
	app.Run(iris.Addr(":8080"))
}
