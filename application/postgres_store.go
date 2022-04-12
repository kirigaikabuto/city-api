package application

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

var applicationQueries = []string{
	`create table if not exists Applications(
		id text,
		address text,
		app_type text,
		message text,
		first_name text,
		last_name text,
		patronymic text,
		phone_number text,
		photo_url text,
		video_url text,
		created_date date,
		longitude double precision,
		latitude double precision,
		primary key(id)
	);`,
}

type applicationStore struct {
	db *sql.DB
}

func NewPostgresApplicationStore(cfg common.PostgresConfig) (Store, error) {
	db, err := common.GetDbConn(common.GetConnString(cfg))
	if err != nil {
		return nil, err
	}
	for _, q := range applicationQueries {
		_, err = db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
	db.SetMaxOpenConns(10)
	store := &applicationStore{db: db}
	return store, nil
}

func (a *applicationStore) Create(model *Application) (*Application, error) {
	model.Id = uuid.New().String()
	result, err := a.db.Exec(
		"INSERT INTO Applications "+
			"(id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, current_date, $11, $12)",
		model.Id, model.Address, model.AppType.ToString(), model.Message, model.FirstName, model.LastName, model.Patronymic,
		model.PhoneNumber, model.PhotoUrl, model.VideoUrl, model.Longitude, model.Latitude,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateApplicationUnknown
	}
	return model, nil
}

func (a *applicationStore) List() ([]Application, error) {
	var objects []Application
	var values []interface{}
	q := "select " +
		"id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude " +
		"from Applications"
	rows, err := a.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Application{}
		appType := ""
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		obj.AppType = ToProblemType(appType)
		objects = append(objects, obj)
	}
	return objects, nil
}

func (a *applicationStore) GetById(id string) (*Application, error) {
	return nil, nil
}

func (a *applicationStore) GetByProblemType(problemType ProblemType) ([]Application, error) {
	return nil, nil
}
