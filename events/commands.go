package events

import "bytes"

type CreateEventCommand struct {
	*Event
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
