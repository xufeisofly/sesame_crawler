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
	builder.RunWith(db.DB).PlaceholderFormat(sq.Dollar).QueryRow().Scan(&id)

	return id
}

func (db TicketDAO) Update(ticketData *Ticket) int {
	columnValueMap := map[string]interface{}{
		"train_no":   ticketData.TrainNo,
		"start_id":   ticketData.StartId,
		"end_id":     ticketData.EndId,
		"start_time": ticketData.StartTime,
		"end_time":   ticketData.EndTime,
		"updated_at": time.Now(),
	}

	builder := sq.Update("tickets").
		SetMap(columnValueMap).
		Where(sq.Eq{"id": ticketData.Id}).
		Suffix("RETURNING \"id\"")

	var id int
	builder.RunWith(db.DB).PlaceholderFormat(sq.Dollar).QueryRow().Scan(&id)

	return id
}

func (db TicketDAO) GetByRoute(startId, endId int, trainNo string) *Ticket {
	builder := sq.
		Select("*").
		From("tickets").
		Where(sq.Eq{"start_id": startId, "end_id": endId, "train_no": trainNo})
	rows, err := builder.RunWith(db.DB).PlaceholderFormat(sq.Dollar).Query()
	if err != nil {
		panic(err)
	}

	var ret interface{}
	if rows.Next() {
		ret = scanTicket(rows)
	}
	defer rows.Close()

	if ret == nil {
		return &Ticket{}
	}

	return ret.(*Ticket)
}
