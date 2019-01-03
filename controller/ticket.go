package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	uri "net/url"
	"sesame/proxy"
	"strings"
	"time"
)

type Ticket struct {
	TrainNo   string
	StartTime string
	EndTime   string
	Duration  string
}

func GetTickets(from_code, to_code, date string) []Ticket {
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
	req.Header.Set("User-Agent", proxy.GetAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	pxy, err := uri.Parse(proxy.ReturnIp())
	timeout := time.Duration(20 * time.Second)
	fmt.Printf("使用代理:%s\n", pxy)

	if err != nil {
		log.Fatal(err)
	}

	var resp *http.Response
	client := &http.Client{
		// Transport: &http.Transport{
		// 	Proxy:           http.ProxyURL(pxy),
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// },
		Timeout: timeout,
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := ioutil.ReadAll(resp.Body)

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
