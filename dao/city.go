package dao

import (
	"database/sql"
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

func scan(rows *sql.Rows) interface{} {
	b := &City{}
	rows.Scan(&b.Id, &b.Name, &b.Code, &b.CreatedAt, &b.UpdatedAt)
	return b
}

func (db CityDAO) Get(id int32) *City {
	builder := sq.Select("*").From("cities").Where(sq.Eq{"id": id})
	rows, _ := builder.RunWith(db).PlaceholderFormat(sq.Dollar).Query()

	var ret interface{}
	if rows.Next() {
		ret = scan(rows)
	}
	defer rows.Close()

	return ret.(*City)
}

func (db CityDAO) GetBy(column string, value interface{}) *City {
	builder := sq.Select("*").From("cities").Where(sq.Eq{column: value.(string)})

	rows, _ := builder.RunWith(db).PlaceholderFormat(sq.Dollar).Query()

	var ret interface{}
	if rows.Next() {
		ret = scan(rows)
	}
	defer rows.Close()

	return ret.(*City)
}
