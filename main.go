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
	// "天津",
	// "南京",
	// "广州",
	// "哈尔滨",
	// "沈阳",
	// "长春",
	// "呼和浩特",
	// "郑州",
	// "济南",
	// "杭州",
	// "上海",
	// "厦门",
	// "成都",
	// "重庆",
	// "拉萨",
	// "乌鲁木齐",
	// "西宁",
	// "昆明",
}

func main() {
	db, _ := sql.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=%s",
			user, dbname, sslmode))

	cityDao := dao.CityDAO{db}
	defer cityDao.Close()

	fromCity := cityDao.GetBy("name", "北京")

	for _, destination := range destinations {
		toCity := cityDao.GetBy("name", destination)
		controller.GetTickets(fromCity.Code, toCity.Code, "2019-01-13")
	}
}
