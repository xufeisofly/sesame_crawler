package dao

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type CityDAO struct {
	*sql.DB
}

type City struct {
	Id        int
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func scanCity(rows *sql.Rows) interface{} {
	b := &City{}
	rows.Scan(&b.Id, &b.Name, &b.Code, &b.CreatedAt, &b.UpdatedAt)
	return b
}

func (db CityDAO) Get(id int32) *City {
	builder := sq.Select("*").From("cities").Where(sq.Eq{"id": id})
	rows, _ := builder.RunWith(db).PlaceholderFormat(sq.Dollar).Query()

	var ret interface{}
	if rows.Next() {
		ret = scanCity(rows)
	}
	defer rows.Close()

	return ret.(*City)
}

func (db CityDAO) GetBy(column string, value interface{}) *City {
	builder := sq.Select("*").From("cities").Where(sq.Eq{column: value.(string)})

	rows, _ := builder.RunWith(db).PlaceholderFormat(sq.Dollar).Query()

	var ret interface{}
	if rows.Next() {
		ret = scanCity(rows)
	}
	defer rows.Close()

	return ret.(*City)
}

func (db CityDAO) MGetByTag(value int) []*City {
	builder := sq.
		Select("cities.*").
		From("cities").
		LeftJoin("city_tag_relations AS r ON r.city_id =  cities.id").
		LeftJoin("tags ON tags.id = r.tag_id").
		Where(sq.Eq{"tags.category": value})

	fmt.Println(db)

	rows, err := builder.RunWith(db).PlaceholderFormat(sq.Dollar).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var ret []*City
	for rows.Next() {
		ret = append(ret, scanCity(rows).(*City))
	}
	return ret
}
