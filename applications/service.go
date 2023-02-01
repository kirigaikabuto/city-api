package applications

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	"github.com/kirigaikabuto/city-api/users"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	twoGisApiKey = ""
	twoGisCity   = ""
)

type Service interface {
	CreateApplication(cmd *CreateApplicationCommand) (*Application, error)
	ListApplications(cmd *ListApplicationsCommand) (*ListApplicationResponse, error)
	UploadApplicationFile(cmd *UploadApplicationFileCommand) (*UploadApplicationFileResponse, error)
	ListApplicationsByType(cmd *ListApplicationsByTypeCommand) ([]Application, error)
	GetApplicationById(cmd *GetApplicationByIdCommand) (*Application, error)
	UpdateApplicationStatus(cmd *UpdateApplicationStatusCommand) (*Application, error)
	ListApplicationsByUserId(cmd *ListApplicationsByUserIdCommand) ([]Application, error)
	UpdateApplication(cmd *UpdateApplicationCommand) (*Application, error)
	RemoveApplication(cmd *RemoveApplicationCommand) error
	ListByAddress(cmd *ListByAddressCommand) ([]Application, error)
	UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*Application, error)

	SearchPlace(cmd *SearchPlaceCommand) ([]Place, error)
}

type service struct {
	appStore   Store
	filesStore file_storage.Store
	userStore  users.UsersStore
	s3         common.S3Uploader
}

func NewApplicationService(appStore Store, s3 common.S3Uploader, usersStore users.UsersStore, fileStore file_storage.Store) Service {
	return &service{appStore: appStore, s3: s3, userStore: usersStore, filesStore: fileStore}
}

func (s *service) CreateApplication(cmd *CreateApplicationCommand) (*Application, error) {
	var appType ProblemType
	if IsProblemTypeExist(cmd.AppType) {
		appType = ToProblemType(cmd.AppType)
	} else {
		return nil, ErrApplicationTypeNotExist
	}
	app := &Application{
		AppType:     appType,
		Message:     cmd.Message,
		FirstName:   cmd.FirstName,
		LastName:    cmd.LastName,
		Patronymic:  cmd.Patronymic,
		PhoneNumber: cmd.PhoneNumber,
		Address:     cmd.Address,
		Longitude:   cmd.Longitude,
		Latitude:    cmd.Latitude,
	}
	if cmd.UserId != "" {
		app.UserId = cmd.UserId
		currentUser, err := s.userStore.Get(cmd.UserId)
		if err != nil {
			return nil, err
		}
		app.FirstName = currentUser.FirstName
		app.LastName = currentUser.LastName
		app.PhoneNumber = currentUser.PhoneNumber
	}
	return s.appStore.Create(app)
}

func (s *service) ListApplications(cmd *ListApplicationsCommand) (*ListApplicationResponse, error) {
	applications, err := s.appStore.List()
	if err != nil {
		return nil, err
	}
	resp := &ListApplicationResponse{Applications: applications}
	return resp, nil
}

func (s *service) SearchPlace(cmd *SearchPlaceCommand) ([]Place, error) {
	twoGisApiKey = viper.GetString("2gis.primary.api_key")
	twoGisCity = viper.GetString("2gis.primary.city")
	clt := &http.Client{}
	basicUrl := "https://catalog.api.2gis.com/3.0/items/geocode?q=%s, %s&fields=items.point&key=%s"
	finalUrl := fmt.Sprintf(basicUrl, twoGisCity, cmd.Name, twoGisApiKey)
	objects, err := clt.Get(finalUrl)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(objects.Body)
	if err != nil {
		return nil, err
	}
	responseFromApi := &MapResponse{}
	err = json.Unmarshal(data, &responseFromApi)
	if err != nil {
		return nil, err
	}
	fmt.Println(responseFromApi)
	var result []Place
	for _, v := range responseFromApi.Result.Items {
		p := Place{
			Name:        v.Name,
			PurposeName: v.PurposeName,
			FullName:    v.FullName,
			Address:     v.AddressName,
			Type:        v.Type,
			Latitude:    v.Point.Lat,
			Longitude:   v.Point.Lon,
		}
		result = append(result, p)
	}

	return result, nil
}

func (s *service) UploadApplicationFile(cmd *UploadApplicationFileCommand) (*UploadApplicationFileResponse, error) {
	if cmd.ContentType == "" {
		return nil, ErrCannotDetectContentType
	}
	if !common.IsImage(cmd.ContentType) && !common.IsVideo(cmd.ContentType) {
		return nil, ErrFileShouldBeOnlyImageOrVideo
	}
	modelUpdate := &ApplicationUpdate{
		Id: cmd.Id,
	}
	fileType := strings.Split(cmd.ContentType, "/")[1]
	fileResponse, err := s.s3.UploadFile(cmd.File.Bytes(), cmd.Id, fileType, cmd.ContentType)
	if err != nil {
		return nil, err
	}
	if common.IsImage(cmd.ContentType) {
		modelUpdate.PhotoUrl = &fileResponse.FileUrl
	} else if common.IsVideo(cmd.ContentType) {
		modelUpdate.VideoUrl = &fileResponse.FileUrl
	} else {
		return nil, ErrInCheckForContentType
	}
	_, err = s.appStore.Update(modelUpdate)
	if err != nil {
		return nil, err
	}
	response := &UploadApplicationFileResponse{FileUrl: fileResponse.FileUrl}
	return response, nil
}

func (s *service) ListApplicationsByType(cmd *ListApplicationsByTypeCommand) ([]Application, error) {
	var appType ProblemType
	if IsProblemTypeExist(cmd.AppType) {
		appType = ToProblemType(cmd.AppType)
	} else {
		return nil, ErrApplicationTypeNotExist
	}
	return s.appStore.GetByProblemType(appType)
}

func (s *service) GetApplicationById(cmd *GetApplicationByIdCommand) (*Application, error) {
	return s.appStore.GetById(cmd.Id)
}

func (s *service) UpdateApplicationStatus(cmd *UpdateApplicationStatusCommand) (*Application, error) {
	model := &ApplicationUpdate{
		Id: cmd.Id,
	}
	if !IsStatusExist(cmd.Status) {
		return nil, ErrApplicationStatusNotExist
	}
	currentStatus := ToStatus(cmd.Status)
	model.AppStatus = &currentStatus
	return s.appStore.Update(model)
}

func (s *service) ListApplicationsByUserId(cmd *ListApplicationsByUserIdCommand) ([]Application, error) {
	return s.appStore.ListApplicationsByUserId(cmd.UserId)
}

func (s *service) UpdateApplication(cmd *UpdateApplicationCommand) (*Application, error) {
	updateModel := &ApplicationUpdate{Id: cmd.Id}
	currentModel, err := s.appStore.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}
	if cmd.FirstName != nil && *cmd.FirstName != currentModel.FirstName {
		updateModel.FirstName = cmd.FirstName
	}
	if cmd.LastName != nil && *cmd.LastName != currentModel.LastName {
		updateModel.LastName = cmd.LastName
	}
	if cmd.Patronymic != nil && *cmd.Patronymic != currentModel.Patronymic {
		updateModel.Patronymic = cmd.Patronymic
	}
	if cmd.PhoneNumber != nil && *cmd.PhoneNumber != currentModel.PhoneNumber {
		updateModel.PhoneNumber = cmd.PhoneNumber
	}
	if cmd.Address != nil && *cmd.Address != currentModel.Address {
		updateModel.Address = cmd.Address
	}
	if cmd.Latitude != nil && *cmd.Latitude != currentModel.Latitude {
		updateModel.Latitude = cmd.Latitude
	}
	if cmd.Longitude != nil && *cmd.Longitude != currentModel.Longitude {
		updateModel.Longitude = cmd.Longitude
	}
	if cmd.AppType != nil && *cmd.AppType != currentModel.AppType {
		updateModel.AppType = cmd.AppType
	}
	if cmd.Message != nil && *cmd.Message != currentModel.Message {
		updateModel.Message = cmd.Message
	}
	return s.appStore.Update(updateModel)
}

func (s *service) RemoveApplication(cmd *RemoveApplicationCommand) error {
	return s.appStore.RemoveApplication(cmd.Id)
}

func (s *service) ListByAddress(cmd *ListByAddressCommand) ([]Application, error) {
	applications := []Application{}
	tempApplications, err := s.appStore.ListByAddress(cmd.Address)
	if err != nil {
		return nil, err
	}
	if cmd.UserId != "" {
		for _, v := range tempApplications {
			if v.UserId == cmd.UserId {
				applications = append(applications, v)
			}
		}
	} else {
		applications = tempApplications
	}
	return applications, nil
}

func (s *service) UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*Application, error) {
	_, err := s.appStore.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}
	for _, obj := range cmd.Files {
		if obj.ContentType == "" {
			return nil, ErrCannotDetectContentType
		}
		if !common.IsImage(obj.ContentType) && !common.IsVideo(obj.ContentType) {
			return nil, ErrFileShouldBeOnlyImageOrVideo
		}
		fileType := strings.Split(obj.ContentType, "/")[1]
		fileResponse, err := s.s3.UploadFile(obj.File.Bytes(), cmd.Id+uuid.New().String(), fileType, obj.ContentType)
		if err != nil {
			return nil, err
		}
		_, err = s.filesStore.Create(&file_storage.FileStorage{
			ObjectId:   cmd.Id,
			ObjectType: file_storage.ApplicationObjType,
			FileUrl:    fileResponse.FileUrl,
		})
		if err != nil {
			return nil, err
		}
	}
	return s.GetApplicationById(&GetApplicationByIdCommand{Id: cmd.Id})
}
