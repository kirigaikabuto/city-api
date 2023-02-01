package events

import (
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	"strings"
)

type Service interface {
	Create(cmd *CreateEventCommand) (*Event, error)
	List(cmd *ListEventCommand) ([]Event, error)
	ListEventByUserId(cmd *ListEventByUserIdCommand) ([]Event, error)
	UploadDocument(cmd *UploadDocumentCommand) (*UploadDocumentResponse, error)
	GetEventById(cmd *GetEventByIdCommand) (*Event, error)
	UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*Event, error)
}

type service struct {
	eventStore Store
	fileStore  file_storage.Store
	s3         common.S3Uploader
}

func NewService(e Store, s3 common.S3Uploader, fileStore file_storage.Store) Service {
	return &service{eventStore: e, s3: s3, fileStore: fileStore}
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

func (s *service) UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*Event, error) {
	_, err := s.eventStore.GetById(cmd.Id)
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
		_, err = s.fileStore.Create(&file_storage.FileStorage{
			ObjectId:   cmd.Id,
			ObjectType: file_storage.EventObjType,
			FileUrl:    fileResponse.FileUrl,
		})
		if err != nil {
			return nil, err
		}
	}
	return s.GetEventById(&GetEventByIdCommand{Id: cmd.Id})
}
