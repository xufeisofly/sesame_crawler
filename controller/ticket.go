package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Ticket struct {
	TrainNo   string
	StartTime string
	EndTime   string
	Duration  string
}

func GetTickets(from_code, to_code, date string) {
	baseUrl := "https://kyfw.12306.cn/otn/leftTicket/queryZ"

	params := []map[string]string{
		map[string]string{"leftTicketDTO.train_date": date},
		map[string]string{"leftTicketDTO.from_station": from_code},
		map[string]string{"leftTicketDTO.to_station": to_code},
		map[string]string{"purpose_codes": "ADULT"},
	}

	queryUrl := "?"
	for _, item := range params {
		for k, v := range item {
			queryUrl += k + "=" + v + "&"
		}
	}

	url := baseUrl + queryUrl
	url = url[0 : len(url)-1]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := ioutil.ReadAll(resp.Body)

	tickets := dumpData(s)
	// for _, ticket := range tickets {
	// 	fmt.Printf("车次: %s | 时长: %s \n", ticket.TrainNo, ticket.Duration)
	// }
	ticket := tickets[0]
	fmt.Printf("车次: %s | 时长: %s \n", ticket.TrainNo, ticket.Duration)

	defer resp.Body.Close()
}

func dumpData(data []byte) []Ticket {
	var jsonResult map[string]interface{}

	json.Unmarshal(data, &jsonResult)
	results := jsonResult["data"].(map[string]interface{})["result"]
	arr := results.([]interface{})

	var tickets []Ticket
	for _, i := range arr {
		str_arr := strings.Split(i.(string), "|")
		ticket := Ticket{
			TrainNo:   str_arr[3],
			StartTime: str_arr[8],
			EndTime:   str_arr[9],
			Duration:  str_arr[10],
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}
