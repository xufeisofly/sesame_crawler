package main

import (
	"fmt"
	"sesame/controller"
	"sesame/dao"
	"time"

	"sesame/db"
)

func main() {
	// Init DB
	db := db.Database
	// Init city dao
	cityDao := dao.CityDAO{db}
	ticketDao := dao.TicketDAO{db}
	defer cityDao.Close()
	defer ticketDao.Close()

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
					fmt.Printf(
						"Updated %s--%s 车次: %s | 时长: %s \n",
						startCity.Name,
						endCity.Name,
						ticket.TrainNo,
						ticket.Duration)
				} else {
					ticketDao.Create(
						startCity.Id,
						endCity.Id,
						ticket.TrainNo,
						ticket.StartTime,
						ticket.EndTime,
						ticket.Duration,
					)
					fmt.Printf(
						"Created %s--%s 车次: %s | 时长: %s \n",
						startCity.Name,
						endCity.Name,
						ticket.TrainNo,
						ticket.Duration)
				}
			}
			time.Sleep(3 * time.Second)
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
