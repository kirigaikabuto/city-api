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

type ListApplicationResponse struct {
	Applications []Application `json:"applications"`
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
	Id          string        `json:"id"`
	File        *bytes.Buffer `json:"file" form:"file"`
	ContentType string        `json:"-"`
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

type UpdateApplicationStatusCommand struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func (cmd *UpdateApplicationStatusCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UpdateApplicationStatus(cmd)
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

type ListApplicationsByUserIdCommand struct {
	UserId string `json:"-"`
}

func (cmd *ListApplicationsByUserIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListApplicationsByUserId(cmd)
}

type UpdateApplicationCommand struct {
	ApplicationUpdate
}

func (cmd *UpdateApplicationCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UpdateApplication(cmd)
}

type RemoveApplicationCommand struct {
	Id string `json:"id"`
}

func (cmd *RemoveApplicationCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(Service).RemoveApplication(cmd)
}

type ListByAddressCommand struct {
	Address string `json:"address"`
	UserId  string `json:"user_id"`
}

func (cmd *ListByAddressCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListByAddress(cmd)
}

type FileObj struct {
	File        *bytes.Buffer
	ContentType string
}

type UploadMultipleFilesCommand struct {
	Id    string `json:"id"`
	Files []FileObj
}

func (cmd *UploadMultipleFilesCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UploadMultipleFiles(cmd)
}
