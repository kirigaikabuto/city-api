package events

import (
	"github.com/kirigaikabuto/city-api/common"
	"strings"
)

type Service interface {
	Create(cmd *CreateEventCommand) (*Event, error)
	List(cmd *ListEventCommand) ([]Event, error)
	ListEventByUserId(cmd *ListEventByUserIdCommand) ([]Event, error)
	UploadDocument(cmd *UploadDocumentCommand) (*UploadDocumentResponse, error)
	GetEventById(cmd *GetEventByIdCommand) (*Event, error)
}

type service struct {
	eventStore Store
	s3         common.S3Uploader
}

func NewService(e Store, s3 common.S3Uploader) Service {
	return &service{eventStore: e, s3: s3}
}

func (s *service) Create(cmd *CreateEventCommand) (*Event, error) {
	return s.eventStore.Create(&cmd.Event)
}

func (s *service) List(cmd *ListEventCommand) ([]Event, error) {
	return s.eventStore.List()
}

func (s *service) ListEventByUserId(cmd *ListEventByUserIdCommand) ([]Event, error) {
	return s.eventStore.ListByUserId(cmd.UserId)
}

func (s *service) UploadDocument(cmd *UploadDocumentCommand) (*UploadDocumentResponse, error) {
	modelUpdate := &EventUpdate{
		Id: cmd.Id,
	}
	fileType := strings.Split(cmd.ContentType, "/")[1]
	fileResponse, err := s.s3.UploadFile(cmd.File.Bytes(), cmd.Id, fileType, cmd.ContentType)
	if err != nil {
		return nil, err
	}
	modelUpdate.DocumentUrl = &fileResponse.FileUrl
	_, err = s.eventStore.Update(modelUpdate)
	if err != nil {
		return nil, err
	}
	response := &UploadDocumentResponse{FileUrl: fileResponse.FileUrl}
	return response, nil
}

func (s *service) GetEventById(cmd *GetEventByIdCommand) (*Event, error) {
	event, err := s.eventStore.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}
	return event, nil
}
