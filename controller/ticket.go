package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	myproxy "sesame/proxy"
	"strings"

	"golang.org/x/net/proxy"
)

type Ticket struct {
	TrainNo   string
	StartTime string
	EndTime   string
	Duration  string
}

func GetTickets(from_code, to_code, date string) []Ticket {
	baseUri := "https://kyfw.12306.cn/otn/leftTicket/queryZ"

	params := []map[string]string{
		map[string]string{"leftTicketDTO.train_date": date},
		map[string]string{"leftTicketDTO.from_station": from_code},
		map[string]string{"leftTicketDTO.to_station": to_code},
		map[string]string{"purpose_codes": "ADULT"},
	}

	queryUri := "?"
	for _, item := range params {
		for k, v := range item {
			queryUri += k + "=" + v + "&"
		}
	}

	uri := baseUri + queryUri
	uri = uri[0 : len(uri)-1]
	uri = "https://www.niltouch.cn"

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", myproxy.GetAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	// pxy, err := url.Parse(myproxy.ReturnIp())

	// if err != nil {
	// 	log.Fatal(err)
	// }

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	var resp *http.Response
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(s))
	tickets := dumpData(s)
	defer resp.Body.Close()

	return tickets
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
