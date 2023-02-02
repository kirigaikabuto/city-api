package events

import "bytes"

type CreateEventCommand struct {
	Event
}

func (cmd *CreateEventCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Create(cmd)
}

type ListEventCommand struct {
}

func (cmd *ListEventCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).List(cmd)
}

type ListEventByUserIdCommand struct {
	UserId string `json:"user_id"`
}

func (cmd *ListEventByUserIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListEventByUserId(cmd)
}

type UploadDocumentCommand struct {
	UserId      string        `json:"-"`
	Id          string        `json:"id"`
	File        *bytes.Buffer `json:"file" form:"file"`
	ContentType string        `json:"-"`
}

func (cmd *UploadDocumentCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UploadDocument(cmd)
}

type UploadDocumentResponse struct {
	FileUrl string `json:"file_url"`
}

type GetEventByIdCommand struct {
	Id string `json:"id"`
}

func (cmd *GetEventByIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetEventById(cmd)
}

type GetEventByIdResponse struct {
	Event
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
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
