package events

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	"log"
	"strconv"
	"strings"
)

var queries = []string{
	`create table if not exists events (
		id text,
		address text,
		description text,
		date date,
		time text,
		organizer_info text,
		document_url text,
		longitude double precision,
		latitude double precision,
		created_date date,
		user_id text,
		primary key(id)
	);`,
}

type store struct {
	db *sql.DB
}

func NewPostgresStore(cfg common.PostgresConfig) (Store, error) {
	db, err := common.GetDbConn(common.GetConnString(cfg))
	if err != nil {
		return nil, err
	}
	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
	db.SetMaxOpenConns(10)
	s := &store{db: db}
	return s, nil
}

func (s *store) Create(model *Event) (*Event, error) {
	model.Id = uuid.New().String()
	result, err := s.db.Exec(
		"INSERT INTO events "+
			"(id, address, description, date, time, organizer_info, document_url, longitude, latitude, created_date, user_id) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, current_date, $10)",
		model.Id, model.Address, model.Description, model.Date, model.Time, model.OrganizerInfo, model.DocumentUrl, model.Longitude,
		model.Latitude, model.UserId,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateEventUnknown
	}
	return model, nil
}

func (s *store) List() ([]Event, error) {
	var objects []Event
	var values []interface{}
	q := "select " +
		"id, address, description, date, time, organizer_info, document_url, longitude, latitude, created_date, user_id " +
		"from events"
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Event{}
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&obj.Description, &obj.Date,
			&obj.Time, &obj.OrganizerInfo,
			&obj.DocumentUrl, &obj.Longitude,
			&obj.Latitude, &obj.CreatedDate, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.Date = strings.Split(obj.Date, "T")[0]
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) ListByUserId(userId string) ([]Event, error) {
	var objects []Event
	var values []interface{}
	q := "select " +
		"id, address, description, date, time, organizer_info, document_url, longitude, latitude, created_date, user_id " +
		"from events where user_id = $1"
	values = append(values, userId)
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Event{}
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&obj.Description, &obj.Date,
			&obj.Time, &obj.OrganizerInfo,
			&obj.DocumentUrl, &obj.Longitude,
			&obj.Latitude, &obj.CreatedDate, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.Date = strings.Split(obj.Date, "T")[0]
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) GetById(id string) (*Event, error) {
	obj := &Event{}
	err := s.db.QueryRow("select id, address, description, "+
		"date, time, organizer_info, document_url, longitude, latitude, created_date, user_id from events where id = $1", id).
		Scan(&obj.Id, &obj.Address,
			&obj.Description, &obj.Date,
			&obj.Time, &obj.OrganizerInfo,
			&obj.DocumentUrl, &obj.Longitude,
			&obj.Latitude, &obj.CreatedDate, &obj.UserId)
	if err == sql.ErrNoRows {
		return nil, ErrEventNotFound
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *store) Update(model *EventUpdate) (*Event, error) {
	q := "update events set "
	parts := []string{}
	values := []interface{}{}
	cnt := 0
	if model.DocumentUrl != nil {
		cnt++
		parts = append(parts, "document_url = $"+strconv.Itoa(cnt))
		values = append(values, model.DocumentUrl)
	}
	if len(parts) <= 0 {
		return nil, ErrNothingToUpdate
	}
	cnt++
	q = q + strings.Join(parts, " , ") + " WHERE id = $" + strconv.Itoa(cnt)
	values = append(values, model.Id)
	result, err := s.db.Exec(q, values...)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrEventNotFound
	}
	return s.GetById(model.Id)
}
