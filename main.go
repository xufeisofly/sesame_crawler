package main

import (
	"database/sql"
	"fmt"
	"sesame/controller"
	"sesame/dao"
	"time"

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

			for _, ticket := range tickets {
				oldTicket := ticketDao.GetByRoute(startCity.Id, endCity.Id, ticket.TrainNo)

				if oldTicket.TrainNo != "" {
					newTicket := dao.Ticket{
						Id:        oldTicket.Id,
						StartId:   startCity.Id,
						EndId:     endCity.Id,
						TrainNo:   ticket.TrainNo,
						StartTime: ticket.StartTime,
						EndTime:   ticket.EndTime,
					}
					ticketDao.Update(&newTicket)
				} else {
					ticketDao.Create(
						startCity.Id,
						endCity.Id,
						ticket.TrainNo,
						ticket.StartTime,
						ticket.EndTime,
						ticket.Duration,
					)
				}

				fmt.Printf("车次: %s | 时长: %s \n", ticket.TrainNo, ticket.Duration)
			}
			time.Sleep(1 * time.Second)
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

// 	ticket := ticketDao.GetByRoute(26, 75, "K2632")
// 	newTicket := dao.Ticket{
// 		Id:        ticket.Id,
// 		StartId:   99,
// 		EndId:     ticket.EndId,
// 		StartTime: "hey",
// 		EndTime:   "hey",
// 		TrainNo:   "KKK",
// 	}
// 	id := ticketDao.Update(&newTicket)

// 	fmt.Println(id)
// }
