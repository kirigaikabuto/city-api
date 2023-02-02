package events

import (
	"github.com/google/uuid"
	"github.com/kirigaikabuto/city-api/common"
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	"github.com/kirigaikabuto/city-api/users"
	"strings"
)

type Service interface {
	Create(cmd *CreateEventCommand) (*Event, error)
	List(cmd *ListEventCommand) ([]GetEventByIdResponse, error)
	ListEventByUserId(cmd *ListEventByUserIdCommand) ([]GetEventByIdResponse, error)
	UploadDocument(cmd *UploadDocumentCommand) (*UploadDocumentResponse, error)
	GetEventById(cmd *GetEventByIdCommand) (*GetEventByIdResponse, error)
	UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*GetEventByIdResponse, error)
}

type service struct {
	eventStore Store
	fileStore  file_storage.Store
	s3         common.S3Uploader
	usersStore users.UsersStore
}

func NewService(e Store, s3 common.S3Uploader, fileStore file_storage.Store, usersStore users.UsersStore) Service {
	return &service{eventStore: e, s3: s3, fileStore: fileStore, usersStore: usersStore}
}

func (s *service) Create(cmd *CreateEventCommand) (*Event, error) {
	return s.eventStore.Create(&cmd.Event)
}

func (s *service) List(cmd *ListEventCommand) ([]GetEventByIdResponse, error) {
	events, err := s.eventStore.List()
	if err != nil {
		return nil, err
	}
	response := []GetEventByIdResponse{}
	for i := range events {
		user, err := s.usersStore.Get(events[i].UserId)
		if err != nil {
			return nil, err
		}
		resp := GetEventByIdResponse{}
		resp.Event = events[i]
		resp.Username = user.Username
		resp.PhotoUrl = user.Avatar
		response = append(response, resp)
	}
	return response, nil
}

func (s *service) ListEventByUserId(cmd *ListEventByUserIdCommand) ([]GetEventByIdResponse, error) {
	events, err := s.eventStore.ListByUserId(cmd.UserId)
	if err != nil {
		return nil, err
	}
	response := []GetEventByIdResponse{}
	for i := range events {
		user, err := s.usersStore.Get(events[i].UserId)
		if err != nil {
			return nil, err
		}
		resp := GetEventByIdResponse{}
		resp.Event = events[i]
		resp.Username = user.Username
		resp.PhotoUrl = user.Avatar
		response = append(response, resp)
	}
	return response, nil
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

func (s *service) GetEventById(cmd *GetEventByIdCommand) (*GetEventByIdResponse, error) {
	event, err := s.eventStore.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}
	user, err := s.usersStore.Get(event.UserId)
	if err != nil {
		return nil, err
	}
	response := &GetEventByIdResponse{}
	response.Event = *event
	response.Username = user.Username
	response.PhotoUrl = user.Avatar
	return response, nil
}

func (s *service) UploadMultipleFiles(cmd *UploadMultipleFilesCommand) (*GetEventByIdResponse, error) {
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
