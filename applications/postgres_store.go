package applications

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	_ "github.com/lib/pq"
	"log"
	"strconv"
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
		status text,
		user_id text,
		primary key(id)
	);`,
}

type applicationStore struct {
	db               *sql.DB
	fileStorageStore file_storage.Store
}

func NewPostgresApplicationStore(cfg common.PostgresConfig, fileStorageStore file_storage.Store) (Store, error) {
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
	store := &applicationStore{db: db, fileStorageStore: fileStorageStore}
	return store, nil
}

func (a *applicationStore) Create(model *Application) (*Application, error) {
	model.AppStatus = StatusWait
	model.Id = uuid.New().String()
	result, err := a.db.Exec(
		"INSERT INTO Applications "+
			"(id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, current_date, $11, $12, $13, $14)",
		model.Id, model.Address, model.AppType.ToString(), model.Message, model.FirstName, model.LastName, model.Patronymic,
		model.PhoneNumber, model.PhotoUrl, model.VideoUrl, model.Longitude, model.Latitude, model.AppStatus.ToString(), model.UserId,
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
	files, err := a.fileStorageStore.ListByObjectId(model.Id)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		model.Files = append(model.Files, file.FileUrl)
	}
	return model, nil
}

func (a *applicationStore) List() ([]Application, error) {
	var objects []Application
	var values []interface{}
	q := "select " +
		"id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id " +
		"from Applications"
	rows, err := a.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Application{}
		appType := ""
		appStatus := ""
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude, &appStatus, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		obj.AppType = ToProblemType(appType)
		obj.AppStatus = ToStatus(appStatus)
		objects = append(objects, obj)
	}
	for i := range objects {
		files, err := a.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}

func (a *applicationStore) GetById(id string) (*Application, error) {
	obj := &Application{}
	appType := ""
	appStatus := ""
	err := a.db.QueryRow("select id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id from Applications where id = $1", id).
		Scan(&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude,
			&appStatus, &obj.UserId,
		)
	if err == sql.ErrNoRows {
		return nil, ErrApplicationNotFound
	} else if err != nil {
		return nil, err
	}
	obj.AppType = ToProblemType(appType)
	obj.AppStatus = ToStatus(appStatus)
	files, err := a.fileStorageStore.ListByObjectId(obj.Id)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		obj.Files = append(obj.Files, file.FileUrl)
	}
	return obj, nil
}

func (a *applicationStore) GetByProblemType(problemType ProblemType) ([]Application, error) {
	var objects []Application
	var values []interface{}
	values = append(values, problemType.ToString())
	q := "select " +
		"id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id " +
		"from Applications where app_type = $1"
	rows, err := a.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Application{}
		appType := ""
		appStatus := ""
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude, &appStatus, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		obj.AppType = ToProblemType(appType)
		obj.AppStatus = ToStatus(appStatus)
		objects = append(objects, obj)
	}
	for i := range objects {
		files, err := a.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}

func (a *applicationStore) Update(model *ApplicationUpdate) (*Application, error) {
	q := "update applications set "
	parts := []string{}
	values := []interface{}{}
	cnt := 0
	if model.AppType != nil {
		cnt++
		parts = append(parts, "app_type = $"+strconv.Itoa(cnt))
		values = append(values, model.AppType.ToString())
	}
	if model.Message != nil {
		cnt++
		parts = append(parts, "message = $"+strconv.Itoa(cnt))
		values = append(values, model.Message)
	}
	if model.FirstName != nil {
		cnt++
		parts = append(parts, "first_name = $"+strconv.Itoa(cnt))
		values = append(values, model.FirstName)
	}
	if model.LastName != nil {
		cnt++
		parts = append(parts, "last_name = $"+strconv.Itoa(cnt))
		values = append(values, model.LastName)
	}
	if model.Patronymic != nil {
		cnt++
		parts = append(parts, "patronymic = $"+strconv.Itoa(cnt))
		values = append(values, model.Patronymic)
	}
	if model.Address != nil {
		cnt++
		parts = append(parts, "address = $"+strconv.Itoa(cnt))
		values = append(values, model.Address)
	}
	if model.Latitude != nil {
		cnt++
		parts = append(parts, "latitude = $"+strconv.Itoa(cnt))
		values = append(values, model.Latitude)
	}
	if model.Longitude != nil {
		cnt++
		parts = append(parts, "longitude = $"+strconv.Itoa(cnt))
		values = append(values, model.Longitude)
	}
	if model.PhotoUrl != nil {
		cnt++
		parts = append(parts, "photo_url = $"+strconv.Itoa(cnt))
		values = append(values, model.PhotoUrl)
	}
	if model.VideoUrl != nil {
		cnt++
		parts = append(parts, "video_url = $"+strconv.Itoa(cnt))
		values = append(values, model.VideoUrl)
	}
	if model.AppStatus != nil {
		cnt++
		parts = append(parts, "status = $"+strconv.Itoa(cnt))
		values = append(values, model.AppStatus.ToString())
	}
	if len(parts) <= 0 {
		return nil, ErrNothingToUpdate
	}
	cnt++
	q = q + strings.Join(parts, " , ") + " WHERE id = $" + strconv.Itoa(cnt)
	values = append(values, model.Id)
	result, err := a.db.Exec(q, values...)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, ErrApplicationNotFound
	}
	return a.GetById(model.Id)
}

func (a *applicationStore) ListApplicationsByUserId(userId string) ([]Application, error) {
	var objects []Application
	var values []interface{}
	values = append(values, userId)
	q := "select " +
		"id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id " +
		"from Applications where user_id = $1"
	rows, err := a.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Application{}
		appType := ""
		appStatus := ""
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude, &appStatus, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		obj.AppType = ToProblemType(appType)
		obj.AppStatus = ToStatus(appStatus)
		objects = append(objects, obj)
	}
	for i := range objects {
		files, err := a.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}

func (a *applicationStore) RemoveApplication(id string) error {
	_, err := a.db.Exec("DELETE FROM Applications WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (a *applicationStore) ListByAddress(address string) ([]Application, error) {
	var objects []Application
	var values []interface{}
	q := "select " +
		"id, address, app_type, message, first_name, last_name, patronymic, phone_number, photo_url, video_url, created_date, longitude, latitude, status, user_id " +
		"from Applications where address like '%" + address + "%'"
	rows, err := a.db.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := Application{}
		appType := ""
		appStatus := ""
		err = rows.Scan(
			&obj.Id, &obj.Address,
			&appType, &obj.Message, &obj.FirstName,
			&obj.LastName, &obj.Patronymic,
			&obj.PhoneNumber, &obj.PhotoUrl,
			&obj.VideoUrl, &obj.CreatedDate,
			&obj.Longitude, &obj.Latitude, &appStatus, &obj.UserId)
		if err != nil {
			return nil, err
		}
		obj.CreatedDate = strings.Split(obj.CreatedDate, "T")[0]
		obj.AppType = ToProblemType(appType)
		obj.AppStatus = ToStatus(appStatus)
		objects = append(objects, obj)
	}
	for i := range objects {
		files, err := a.fileStorageStore.ListByObjectId(objects[i].Id)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			objects[i].Files = append(objects[i].Files, file.FileUrl)
		}
	}
	return objects, nil
}
