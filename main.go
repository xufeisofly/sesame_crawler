package main

import (
	"log"
	"math/rand"
	"os"
	"sesame/controller"
	"sesame/dao"
	"time"

	"sesame/db"
)

func sync() {
	// Init DB
	db := db.Database
	// Init city dao
	cityDao := dao.CityDAO{db}
	ticketDao := dao.TicketDAO{db}
	defer cityDao.Close()
	defer ticketDao.Close()

	// Get All Stations By Tag
	cities := cityDao.MGetByTag(1)

	f, err := os.OpenFile("application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	for x, startCity := range cities {
		for y, endCity := range cities {
			// 显示进度
			log.Printf("%v/%v ============ \n", x*len(cities)+y, len(cities)*len(cities))
			if startCity.Name == endCity.Name {
				continue
			}

			if controller.HasSynced(startCity.Name, endCity.Name) {
				log.Printf("%s-%s 已同步", startCity.Name, endCity.Name)
				continue
			}

			tickets := controller.GetTickets(startCity.Name, endCity.Name, "2019-02-01")
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
					log.Printf(
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
					log.Printf(
						"Created %s--%s 车次: %s | 时长: %s \n",
						startCity.Name,
						endCity.Name,
						ticket.TrainNo,
						ticket.Duration)
				}
			}
			// 延时
			secCount := 1 + rand.Intn(1)
			log.Printf("delay %v seconds", secCount)
			time.Sleep(time.Duration(secCount) * time.Second)
		}
	}
	controller.ClearSynced()
	log.Println("同步完毕！清空sync记录")
}

func main() {
	sync()
}
