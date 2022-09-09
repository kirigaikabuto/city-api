package user_events

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	_ "github.com/lib/pq"
	"log"
)

var marketplaceAppRepoQueries = []string{
	`CREATE TABLE IF NOT EXISTS user_events(
		id TEXT,
		user_id TEXT,
		event_id TEXT,
		created_date TEXT,
		PRIMARY KEY(id)
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
	for _, q := range marketplaceAppRepoQueries {
		_, err = db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
	db.SetMaxOpenConns(10)
	st := &store{db: db}
	return st, nil
}

func (s *store) Create(obj *UserEvent) (*UserEvent, error) {
	obj.Id = uuid.New().String()
	result, err := s.db.Exec(
		"INSERT INTO user_events "+
			"(id, user_id, event_id, created_date) "+
			"VALUES ($1, $2, $3, current_date)",
		obj.Id, obj.UserId, obj.EventId,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateUserEventUnknown
	}
	return s.GetUserEventById(obj.Id)
}

func (s *store) ListByEventId(id string) ([]UserEvent, error) {
	var objects []UserEvent
	var values []interface{}
	q := "select " +
		"id, user_id, event_id, created_date " +
		"from user_events where event_id=$1"
	values = append(values, id)
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := UserEvent{}
		err = rows.Scan(
			&obj.Id, &obj.UserId,
			&obj.EventId, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) ListByUserId(id string) ([]UserEvent, error) {
	var objects []UserEvent
	var values []interface{}
	q := "select " +
		"id, user_id, event_id, created_date " +
		"from user_events where user_id=$1"
	values = append(values, id)
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := UserEvent{}
		err = rows.Scan(
			&obj.Id, &obj.UserId,
			&obj.EventId, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) List() ([]UserEvent, error) {
	var objects []UserEvent
	var values []interface{}
	q := "select " +
		"id, user_id, event_id, created_date " +
		"from user_events"
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := UserEvent{}
		err = rows.Scan(
			&obj.Id, &obj.UserId,
			&obj.EventId, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) GetUserEventById(id string) (*UserEvent, error) {
	obj := &UserEvent{}
	err := s.db.QueryRow("select id, user_id, event_id, created_date from user_events where id = $1", id).
		Scan(&obj.Id, &obj.UserId,
			&obj.EventId, &obj.CreatedDate)
	if err == sql.ErrNoRows {
		return nil, ErrCreateUserEventUnknown
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}
