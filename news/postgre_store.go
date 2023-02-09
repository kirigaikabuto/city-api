package news

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)

var marketplaceAppRepoQueries = []string{
	``,
}

type newsStore struct {
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
	store := &newsStore{db: db}
	return store, nil
}

func (u *newsStore) Create(obj *News) (*News, error) {
	obj.Id = uuid.New().String()
	result, err := u.db.Exec(
		"INSERT INTO applications_news "+
			"(id, title, small_description, description, photo_url, author_id, created_date, created_at, modified_at) "+
			"VALUES ($1, $2, $3, $4, $5 , $6, current_date, current_date, current_date)",
		obj.Id, obj.Title, obj.SmallDescription, obj.Description, obj.PhotoUrl, obj.AuthorId,
	)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrCreateNewsUnknown
	}
	return u.GetById(obj.Id)
}

func (u *newsStore) Update(obj *UpdateNews) (*News, error) {
	q := "update applications_news set "
	parts := []string{}
	values := []interface{}{}
	cnt := 0
	if obj.Title != nil {
		cnt++
		parts = append(parts, "title = $"+strconv.Itoa(cnt))
		values = append(values, obj.Title)
	}
	if obj.SmallDescription != nil {
		cnt++
		parts = append(parts, "small_description = $"+strconv.Itoa(cnt))
		values = append(values, obj.SmallDescription)
	}
	if obj.Description != nil {
		cnt++
		parts = append(parts, "description = $"+strconv.Itoa(cnt))
		values = append(values, obj.Description)
	}
	if obj.PhotoUrl != nil {
		cnt++
		parts = append(parts, "photo_url = $"+strconv.Itoa(cnt))
		values = append(values, obj.PhotoUrl)
	}
	if len(parts) <= 0 {
		return nil, ErrNothingToUpdate
	}
	cnt++
	q = q + strings.Join(parts, " , ") + " WHERE id = $" + strconv.Itoa(cnt)
	values = append(values, obj.Id)
	result, err := u.db.Exec(q, values...)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrNewsNotFound
	}
	return u.GetById(obj.Id)
}

func (u *newsStore) List() ([]News, error) {
	var objects []News
	var values []interface{}
	q := "select " +
		"id, title, small_description, description, photo_url, author_id, created_date " +
		"from applications_news"
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := News{}
		err = rows.Scan(
			&obj.Id, &obj.Title,
			&obj.SmallDescription, &obj.Description,
			&obj.PhotoUrl, &obj.AuthorId, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (u *newsStore) GetById(id string) (*News, error) {
	obj := &News{}
	err := u.db.QueryRow("select id, title, small_description,"+
		" description, photo_url, author_id, created_date from applications_news where id = $1", id).
		Scan(&obj.Id, &obj.Title,
			&obj.SmallDescription, &obj.Description,
			&obj.PhotoUrl, &obj.AuthorId, &obj.CreatedDate)
	if err == sql.ErrNoRows {
		return nil, ErrNewsNotFound
	} else if err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *newsStore) GetByAuthorId(authorId string) ([]News, error) {
	var objects []News
	var values []interface{}
	q := "select " +
		"id, title, small_description, description, photo_url, author_id, created_date " +
		"from applications_news where author_id=$1"
	values = append(values, authorId)
	rows, err := u.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := News{}
		err = rows.Scan(
			&obj.Id, &obj.Title,
			&obj.SmallDescription, &obj.Description,
			&obj.PhotoUrl, &obj.AuthorId, &obj.CreatedDate)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}
