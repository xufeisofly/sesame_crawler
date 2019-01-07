package main

import (
	"database/sql"
	"fmt"
	"sesame/controller"
	"sesame/dao"

	_ "github.com/lib/pq"
)

const (
	user    = "norris"
	dbname  = "sesame_development"
	sslmode = "disable"
)

var destinations = []interface{}{
	"石家庄",
	"天津",
	"南京",
	"广州",
	"哈尔滨",
	"沈阳",
	"长春",
	"呼和浩特",
	"郑州",
	"济南",
	"杭州",
	"上海",
	"厦门",
	"成都",
	"重庆",
	"拉萨",
	"乌鲁木齐",
	"西宁",
	"昆明",
}

func main() {
	// Init DB
	db, _ := sql.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=%s",
			user, dbname, sslmode))

	// Init city dao
	cityDao := dao.CityDAO{db}
	ticketDao := dao.TicketDAO{db}
	defer cityDao.Close()

	// Get All Stations By Tag
	cities := cityDao.MGetByTag(1)

	for _, startCity := range cities {
		for _, endCity := range cities {
			if startCity == endCity {
				continue
			}

			tickets := controller.GetTickets(startCity.Code, endCity.Code, "2019-01-13")
			if len(tickets) == 0 {
				continue
			}

			ticket := tickets[0]
			ticketDao.Create(
				startCity.Id,
				endCity.Id,
				ticket.TrainNo,
				ticket.StartTime,
				ticket.EndTime,
				ticket.Duration,
			)
			fmt.Printf("车次: %s | 时长: %s \n", ticket.TrainNo, ticket.Duration)
		}
	}
}

// func main() {
// 	db, _ := sql.Open(
// 		"postgres",
// 		fmt.Sprintf(
// 			"user=%s dbname=%s sslmode=%s",
// 			user, dbname, sslmode))

// 	ticketDao := dao.TicketDAO{db}
// 	defer ticketDao.Close()
// 	id := ticketDao.Create(1, 69, "G85", "08:00", "09:00", "01:00")
// 	fmt.Println(id)
// }
