package main

import (
	"database/sql"
	"fmt"
	"sesame/dao"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// app := iris.Default()

	db, _ := sql.Open("mysql", "norris@/sesame")
	cityDao := dao.CityDAO{db}
	defer cityDao.Close()

	city := cityDao.Get(1)
	fmt.Println(city.Name)

	// db, err := sql.Open("mysql", "norris@/sesame")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// app.Get("/ping", func(ctx iris.Context) {
	// 	ctx.JSON(iris.Map{
	// 		"message": "pong",
	// 	})
	// })

	// app.Get("/tickets", func(ctx iris.Context) {
	// 	from := ctx.URLParam("from")
	// 	to := ctx.URLParam("to")

	// 	app.Logger().Infof("from: %v, to: %v", from, to)
	// })

	// app.Run(iris.Addr(":8080"))
}
