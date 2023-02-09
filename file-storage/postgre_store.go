package file_storage

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	"log"
	"strings"
)

var queries = []string{
	``,
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

func (s *store) Create(model *FileStorage) (*FileStorage, error) {
	model.Id = uuid.New().String()
	result, err := s.db.Exec(
		"INSERT INTO applications_filestorage "+
			"(id, object_type, object_id, file_url, created_date,created_at, modified_at) "+
			"VALUES ($1, $2, $3, $4, current_date, current_date, current_date)",
		model.Id, model.ObjectType, model.ObjectId, model.FileUrl)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateUnknown
	}
	return model, nil
}

func (s *store) ListByObjectType(objType string) ([]FileStorage, error) {
	var objects []FileStorage
	var values []interface{}
	q := "select " +
		"id, object_type, object_id, file_url, created_date " +
		"from applications_filestorage where object_type = $1"
	values = append(values, objType)
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := FileStorage{}
		err = rows.Scan(
			&obj.Id, &obj.ObjectType,
			&obj.ObjectId, &obj.FileUrl,
			&obj.CreatedDate,
		)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) ListByObjectId(objId string) ([]FileStorage, error) {
	var objects []FileStorage
	var values []interface{}
	q := "select " +
		"id, object_type, object_id, file_url, created_date " +
		"from applications_filestorage where object_id = $1"
	values = append(values, objId)
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := FileStorage{}
		err = rows.Scan(
			&obj.Id, &obj.ObjectType,
			&obj.ObjectId, &obj.FileUrl,
			&obj.CreatedDate,
		)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *store) List() ([]FileStorage, error) {
	var objects []FileStorage
	var values []interface{}
	q := "select " +
		"id, object_type, object_id, file_url, created_date " +
		"from applications_filestorage"
	rows, err := s.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := FileStorage{}
		err = rows.Scan(
			&obj.Id, &obj.ObjectType,
			&obj.ObjectId, &obj.FileUrl,
			&obj.CreatedDate,
		)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		objects = append(objects, obj)
	}
	return objects, nil
}
