package applications

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

var (
	twoGisApiKey = ""
	twoGisCity   = ""
)

type Service interface {
	CreateApplication(cmd *CreateApplicationCommand) (*Application, error)
	ListApplications(cmd *ListApplicationsCommand) ([]Application, error)
	UploadApplicationFile(cmd *UploadApplicationFileCommand) (*UploadApplicationFileResponse, error)
	ListApplicationsByType(cmd *ListApplicationsByTypeCommand) ([]Application, error)
	GetApplicationById(cmd *GetApplicationByIdCommand) (*Application, error)

	SearchPlace(cmd *SearchPlaceCommand) ([]Place, error)
}

type service struct {
	appStore Store
	s3       common.S3Uploader
}

func NewApplicationService(appStore Store, s3 common.S3Uploader) Service {
	return &service{appStore: appStore, s3: s3}
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
	return s.appStore.Create(app)
}

func (s *service) ListApplications(cmd *ListApplicationsCommand) ([]Application, error) {
	return s.appStore.List()
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
	fileResponse, err := s.s3.UploadFile(cmd.File.Bytes(), uuid.New().String(), "png")
	if err != nil {
		return nil, err
	}
	modelUpdate := &ApplicationUpdate{
		Id:       cmd.Id,
		PhotoUrl: &fileResponse.FileUrl,
		VideoUrl: nil,
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
