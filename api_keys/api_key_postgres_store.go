package api_keys

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	"log"
	"strings"
)

var queries = []string{
	`create table if not exists api_keys (
		id text,
		key text,
		created_date date,
		primary key(id)
	);`,
}

type store struct {
	db *sql.DB
}

func NewPostgresStore(cfg common.PostgresConfig) (ApiKeyStore, error) {
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

func (s *store) Create() (*ApiKey, error) {
	model := &ApiKey{}
	model.Id = uuid.New().String()
	model.Key = uuid.New().String()
	result, err := s.db.Exec(
		"INSERT INTO api_keys "+
			"(id, key, created_date) "+
			"VALUES ($1, $2, current_date)",
		model.Id, model.Key,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateApiKeyUnknown
	}
	return model, nil
}

func (s *store) Delete(id string) error {
	return nil
}

func (s *store) GetByKey(key string) (*ApiKey, error) {
	obj := &ApiKey{}
	err := s.db.QueryRow("select id, key, created_date from api_keys where key = $1", key).
		Scan(&obj.Id, &obj.Key, &obj.CreatedDate)
	if err == sql.ErrNoRows {
		return nil, ErrApiKeyNotFound
	} else if err != nil {
		return nil, err
	}
	obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
	return obj, nil
}

func (s *store) List() ([]ApiKey, error) {
	var objects []ApiKey
	var values []interface{}
	q := "select " +
		"id, key, created_date " +
		"from api_keys"
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := ApiKey{}
		err = rows.Scan(
			&obj.Id, &obj.Key,
			&obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}
