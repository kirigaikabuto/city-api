package comments

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	_ "github.com/lib/pq"
	"log"
)

var marketplaceAppRepoQueries = []string{
	`CREATE TABLE IF NOT EXISTS comments(
		id TEXT,
		message TEXT,
		user_id TEXT,
		obj_id TEXT,
		obj_type TEXT,
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

func (u *usersStore) Create(obj *Comment) (*Comment, error) {
	obj.Id = uuid.New().String()
	result, err := u.db.Exec(
		"INSERT INTO comments "+
			"(id, message, user_id, obj_id, obj_type, created_date) "+
			"VALUES ($1, $2, $3, $4, $5, current_date)",
		obj.Id, obj.Message, obj.UserId, obj.ObjId, obj.ObjType.ToString(),
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateCommentUnknown
	}
	return u.GetById(obj.Id)
}

func (u *usersStore) List() ([]Comment, error) {
	var objects []Comment
	var values []interface{}
	q := "select " +
		"id, message, user_id, obj_id, obj_type, created_date " +
		"from comments"
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Comment{}
		err = rows.Scan(
			&obj.Id, &obj.Message, &obj.UserId, &obj.ObjId, &obj.ObjType, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (u *usersStore) GetById(id string) (*Comment, error) {
	obj := &Comment{}
	err := u.db.QueryRow("select id, message, user_id, obj_id, obj_type, created_date from comments where id = $1", id).
		Scan(&obj.Id, &obj.Message, &obj.UserId, &obj.ObjId, &obj.ObjType, &obj.CreatedDate)
	if err == sql.ErrNoRows {
		return nil, ErrCommentNotFound
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *usersStore) GetByObjType(objType ObjType) ([]Comment, error) {
	var objects []Comment
	var values []interface{}
	q := "select " +
		"id, message, user_id, obj_id, obj_type, created_date " +
		"from comments where obj_type = $1"
	values = append(values, objType)
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Comment{}
		err = rows.Scan(
			&obj.Id, &obj.Message, &obj.UserId, &obj.ObjId, &obj.ObjType, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (u *usersStore) GetByObjId(objId string) ([]Comment, error) {
	var objects []Comment
	var values []interface{}
	q := "select " +
		"id, message, user_id, obj_id, obj_type, created_date " +
		"from comments where obj_id = $1"
	values = append(values, objId)
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Comment{}
		err = rows.Scan(
			&obj.Id, &obj.Message, &obj.UserId, &obj.ObjId, &obj.ObjType, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}
