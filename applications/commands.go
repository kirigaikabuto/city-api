package applications

import "bytes"

type CreateApplicationCommand struct {
	Address     string  `json:"address"`
	AppType     string  `json:"app_type"`
	Message     string  `json:"message"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Patronymic  string  `json:"patronymic"`
	PhoneNumber string  `json:"phone_number"`
	PhotoUrl    string  `json:"photo_url"`
	VideoUrl    string  `json:"video_url"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	UserId      string  `json:"-"`
}

func (cmd *CreateApplicationCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateApplication(cmd)
}

type ListApplicationsCommand struct {
}

func (cmd *ListApplicationsCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListApplications(cmd)
}

type SearchPlaceCommand struct {
	Name string `json:"name"`
}

func (cmd *SearchPlaceCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).SearchPlace(cmd)
}

type UploadApplicationFileCommand struct {
	Id   string        `json:"id"`
	File *bytes.Buffer `json:"file" form:"file"`
}

func (cmd *UploadApplicationFileCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UploadApplicationFile(cmd)
}

type UploadApplicationFileResponse struct {
	FileUrl string `json:"file_url"`
}

type ListApplicationsByTypeCommand struct {
	AppType string `json:"app_type"`
}

func (cmd *ListApplicationsByTypeCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListApplicationsByType(cmd)
}

type GetApplicationByIdCommand struct {
	Id string `json:"id"`
}

func (cmd *GetApplicationByIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetApplicationById(cmd)
}

type Place struct {
	Name        string  `json:"name"`
	PurposeName string  `json:"purpose_name"`
	FullName    string  `json:"full_name"`
	Address     string  `json:"address"`
	Type        string  `json:"type"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type MapResponse struct {
	Result MapResponseResult `json:"result"`
}

type MapResponseResult struct {
	Items []Item `json:"items"`
}

type Item struct {
	AddressName string `json:"address_name"`
	FullName    string `json:"full_name"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Point       Point  `json:"point"`
	PurposeName string `json:"purpose_name"`
	Type        string `json:"type"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
