package application

import (
	"encoding/json"
	"fmt"
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

	SearchPlace(cmd *SearchPlaceCommand) ([]Place, error)
}

type service struct {
	appStore Store
}

func NewApplicationService(appStore Store) Service {
	return &service{appStore: appStore}
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
	responseFromApi := &TwoGisResponse{}
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
