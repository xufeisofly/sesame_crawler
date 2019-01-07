package dao

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type TicketDAO struct {
	*sql.DB
}

type Ticket struct {
	Id        int
	TrainNo   string
	StartId   int
	EndId     int
	StartTime string
	EndTime   string
	Duration  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func scanTicket(rows *sql.Rows) interface{} {
	b := &Ticket{}
	rows.Scan(&b.Id, &b.TrainNo, &b.StartId, &b.EndId,
		&b.StartTime, &b.EndTime, &b.Duration, &b.CreatedAt, &b.UpdatedAt)
	return b
}

func (db TicketDAO) Create(startId int, endId int, trainNo string,
	startTime string, endTime string, duration string) int {
	builder := sq.Insert("tickets").
		Columns("start_id", "end_id", "train_no", "start_time", "end_time", "duration", "created_at", "updated_at").
		Values(startId, endId, trainNo, startTime, endTime, duration, time.Now(), time.Now()).
		Suffix("RETURNING \"id\"")

	var id int
	builder.RunWith(db).PlaceholderFormat(sq.Dollar).QueryRow().Scan(&id)

	return id
}
