package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	myproxy "sesame/proxy"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

type Ticket struct {
	TrainNo   string
	StartTime string
	EndTime   string
	Duration  string
}

func GetTickets(from_code, to_code, date string) []Ticket {
	baseUrl := "https://train.qunar.com/dict/open/s2s.do"
	curTime := time.Now().Unix() * 1000
	params := []map[string]string{
		map[string]string{"dptStation": "北京"},
		map[string]string{"arrStation": "上海"},
		map[string]string{"date": date},
		map[string]string{"type": "normal"},
		map[string]string{"user": "neibu"},
		map[string]string{"source": "site"},
		map[string]string{"start": "1"},
		map[string]string{"num": "500"},
		map[string]string{"sort": "3"},
		map[string]string{"_": strconv.FormatInt(curTime, 10)},
	}

	queryUrl := "?"
	for _, item := range params {
		for k, v := range item {
			queryUrl += k + "=" + v + "&"
		}
	}

	url := baseUrl + queryUrl
	url = url[0 : len(url)-1]
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", myproxy.GetAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

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
	tickets := dumpData(s)
	defer resp.Body.Close()

	return tickets
}

func dumpData(data []byte) []Ticket {
	var jsonResult map[string]interface{}

	json.Unmarshal(data, &jsonResult)
	results := jsonResult["data"].(map[string]interface{})["s2sBeanList"]

	var tickets []Ticket
	for _, train := range results.([]interface{}) {
		train := train.(map[string]interface{})
		ticket := Ticket{
			TrainNo:   train["trainNo"].(string),
			StartTime: train["dptTime"].(string),
			EndTime:   train["arrTime"].(string),
			Duration:  train["arrTime"].(string),
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}
