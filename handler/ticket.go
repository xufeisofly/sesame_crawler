package handler

import "github.com/kataras/iris"

func TicketList(ctx iris.Context) {
	from := ctx.URLParam("from")
	to := ctx.URLParam("to")

	ctx.JSON(iris.Map{"from": from, "to": to})
}
