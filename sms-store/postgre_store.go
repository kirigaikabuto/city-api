package sms_store

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	"log"
)

var smsStoreRepoQueries = []string{
	`CREATE TABLE IF NOT EXISTS sms_store(
		id TEXT,
		title TEXT,
		type TEXT,
		from_phone TEXT,
		to_phone TEXT,
		body TEXT,
		created_at TEXT,
		PRIMARY KEY(id)
	);`,
}

type smsStore struct {
	db *sql.DB
}

func NewPostgresStore(cfg common.PostgresConfig) (Store, error) {
	db, err := common.GetDbConn(common.GetConnString(cfg))
	if err != nil {
		return nil, err
	}
	for _, q := range smsStoreRepoQueries {
		_, err = db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
	db.SetMaxOpenConns(10)
	store := &smsStore{db: db}
	return store, nil
}

func (u *smsStore) Create(obj *SmsCode) (*SmsCode, error) {
	obj.Id = uuid.New().String()
	result, err := u.db.Exec(
		"INSERT INTO sms_store "+
			"(id, title, type, from_phone, to_phone, body, created_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, current_date)",
		obj.Id, obj.Title, obj.Type, obj.From, obj.To, obj.Body,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateSmsUnknown
	}
	return u.GetById(obj.Id)
}

func (u *smsStore) GetById(id string) (*SmsCode, error) {
	obj := &SmsCode{}
	err := u.db.QueryRow("select id, title, type, from_phone, to_phone, body, created_at from sms_store where id = $1", id).
		Scan(&obj.Id, &obj.Title,
			&obj.Type, &obj.From,
			&obj.To, &obj.Body,
			&obj.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrSmsNotFound
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}
