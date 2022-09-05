package feedback

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	_ "github.com/lib/pq"
	"log"
)

var marketplaceAppRepoQueries = []string{
	`CREATE TABLE IF NOT EXISTS feedback(
		id TEXT,
		message TEXT,
		full_name TEXT,
		phone_number TEXT,
		created_date TEXT,
		PRIMARY KEY(id)
	);`,
}

type usersStore struct {
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
	store := &usersStore{db: db}
	return store, nil
}

func (u *usersStore) Create(obj *Feedback) (*Feedback, error) {
	obj.Id = uuid.New().String()
	result, err := u.db.Exec(
		"INSERT INTO feedback "+
			"(id, message, full_name, phone_number, created_date) "+
			"VALUES ($1, $2, $3, $4, current_date)",
		obj.Id, obj.Message, obj.FullName, obj.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateFeedbackUnknown
	}
	return u.GetById(obj.Id)
}

func (u *usersStore) List() ([]Feedback, error) {
	var objects []Feedback
	var values []interface{}
	q := "select " +
		"id, message, full_name, phone_number, created_date " +
		"from feedback"
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Feedback{}
		err = rows.Scan(
			&obj.Id, &obj.Message,
			&obj.FullName, &obj.PhoneNumber, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (u *usersStore) GetById(id string) (*Feedback, error) {
	obj := &Feedback{}
	err := u.db.QueryRow("select id, message, full_name, phone_number, created_date from feedback where id = $1", id).
		Scan(&obj.Id, &obj.Message,
			&obj.FullName, &obj.PhoneNumber, &obj.CreatedDate)
	if err == sql.ErrNoRows {
		return nil, ErrFeedbackNotFound
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}
