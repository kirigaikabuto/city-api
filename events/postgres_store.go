package events

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	"log"
	"strconv"
	"strings"
)

var queries = []string{
	``,
}

type store struct {
	db               *sql.DB
	fileStorageStore file_storage.Store
}

func NewPostgresStore(cfg common.PostgresConfig, fileStorageStore file_storage.Store) (Store, error) {
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
	s := &store{db: db, fileStorageStore: fileStorageStore}
	return s, nil
}

func (s *store) Create(model *Event) (*Event, error) {
	model.Id = uuid.New().String()
	result, err := s.db.Exec(
		"INSERT INTO applications_event "+
			"(id, address, description, date, time, organizer_info, document_url, longitude, latitude, user_id, created_date, created_at, modified_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, current_date,current_date,current_date)",
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
	files, err := s.fileStorageStore.ListByObjectId(model.Id)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		model.Files = append(model.Files, file.FileUrl)
	}
	return model, nil
}

func (s *store) List() ([]Event, error) {
	var objects []Event
	var values []interface{}
	q := "select " +
		"id, address, description, date, time, organizer_info, document_url, longitude, latitude, created_date, user_id " +
		"from applications_event"
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
	for i := range objects {
		files, err := s.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}

func (s *store) ListByUserId(userId string) ([]Event, error) {
	var objects []Event
	var values []interface{}
	q := "select " +
		"id, address, description, date, time, organizer_info, document_url, longitude, latitude, created_date, user_id " +
		"from applications_event where user_id = $1"
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
	for i := range objects {
		files, err := s.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}

func (s *store) GetById(id string) (*Event, error) {
	obj := &Event{}
	err := s.db.QueryRow("select id, address, description, "+
		"date, time, organizer_info, document_url, longitude, latitude, created_date, user_id from applications_event where id = $1", id).
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
	files, err := s.fileStorageStore.ListByObjectId(obj.Id)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		obj.Files = append(obj.Files, file.FileUrl)
	}
	return obj, nil
}

func (s *store) Update(model *EventUpdate) (*Event, error) {
	q := "update applications_event set "
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
