package proxypool

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomodule/redigo/redis"
)

const PAGE int = 40

var (
	xici string = "http://www.xicidaili.com/wn/"
)

func GetIp(ip string) {

	var count int

	for i := 1; i <= PAGE; i++ {
		response := GetRep(xici+strconv.Itoa(i), ip)

		if response.StatusCode == 200 {
			dom, err := goquery.NewDocumentFromResponse(response)
			if err != nil {
				log.Fatalf("失败原因", response.StatusCode)
			}
			dom.Find("#ip_list tbody tr").Each(func(i int, context *goquery.Selection) {
				// 地址
				ip := context.Find("td").Eq(1).Text()
				// 端口
				port := context.Find("td").Eq(2).Text()

				ip = ip + ":" + port
				saveRedis(ip)
				fmt.Println("获得新代理IP:", ip)
				count++
			})
		}
	}
}

func GetRep(urll string, ip string) *http.Response {
	request, _ := http.NewRequest("GET", urll, nil)
	request.Header.Set("User-Agent", GetAgent())
	request.Header.Set("Accept", "text/html.applicaton/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	proxy, err := url.Parse(ip)

	timeout := time.Duration(20 * time.Second)
	log.Printf("使用代理:%s\n", proxy)
	client := &http.Client{}
	if ip != "local" {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout: timeout,
		}
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		log.Printf("line-99:遇到了错误-并切换ip %s\n", err)
	}

	return response
}

// 随机返回一个User-Agent
func GetAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}

func saveRedis(ip string) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("Connect to redis error", err)
		return
	}

	defer c.Close()
	// 将ip:port 存入set 方便返回随机的ip
	_, err = c.Do("SADD", "IpPool", ip)
	if err != nil {
		log.Fatalf("err:%s", err)
		os.Exit(1)
	}
}

func RemoveRedis(ip string) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("Connect to redis error", err)
		return
	}

	defer c.Close()
	// 将ip:port 存入set 方便返回随机的ip
	_, err = c.Do("SREM", "IpPool", ip)
	if err != nil {
		log.Fatalf("err:%s", err)
		os.Exit(1)
	}
}

func ReturnIp() string {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("Connect to redis error", err)
		os.Exit(0)
	}
	defer c.Close()

	ips, _ := redis.Values(c.Do("SMEMBERS", "IpPool"))

	log.Printf("现有ips %v 个\n", len(ips))
	if len(ips) <= 3 {
		GetIp("local")
	}

	res, err := redis.String(c.Do("SRANDMEMBER", "IpPool"))
	if err != nil {
		log.Println("Random get ip err", err)
	}

	return res
}
