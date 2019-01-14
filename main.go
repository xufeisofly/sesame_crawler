package main

import (
	"fmt"
	"math/rand"
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

			tickets := controller.GetTickets(startCity.Name, endCity.Name, "2019-01-20")
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
			secCount := 10 + rand.Intn(10)
			fmt.Printf("delay %v seconds", secCount)
			time.Sleep(time.Duration(secCount) * time.Second)
		}
	}
}
