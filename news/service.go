package news

import (
	"github.com/kirigaikabuto/city-api/common"
	"strings"
)

type Service interface {
	CreateNews(cmd *CreateNewsCommand) (*News, error)
	ListNews(cmd *ListNewsCommand) ([]News, error)
	UpdateNews(cmd *UpdateNewsCommand) (*News, error)
	GetNewsById(cmd *GetNewsByIdCommand) (*News, error)
	GetNewsByAuthorId(cmd *GetNewsByAuthorId) ([]News, error)
	UploadPhoto(cmd *UploadPhotoCommand) (*UploadPhotoResponse, error)
}

type service struct {
	s3Uploader common.S3Uploader
	store      Store
}

func NewService(s3 common.S3Uploader, store Store) Service {
	return &service{s3Uploader: s3, store: store}
}

func (s *service) CreateNews(cmd *CreateNewsCommand) (*News, error) {
	return s.store.Create(&cmd.News)
}

func (s *service) ListNews(cmd *ListNewsCommand) ([]News, error) {
	return s.store.List()
}

func (s *service) UpdateNews(cmd *UpdateNewsCommand) (*News, error) {
	upd := &UpdateNews{Id: cmd.Id}
	oldNews, err := s.store.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}
	if *cmd.PhotoUrl != "" && oldNews.PhotoUrl != *cmd.PhotoUrl {
		upd.PhotoUrl = cmd.PhotoUrl
	}
	if *cmd.Title != "" && oldNews.Title != *cmd.Title {
		upd.Title = cmd.Title
	}
	if *cmd.SmallDescription != "" && oldNews.SmallDescription != *cmd.SmallDescription {
		upd.SmallDescription = cmd.SmallDescription
	}
	if *cmd.Description != "" && oldNews.Description != *cmd.Description {
		upd.Description = cmd.Description
	}
	if *cmd.AuthorId != "" && oldNews.AuthorId != *cmd.AuthorId {
		upd.AuthorId = cmd.AuthorId
	}
	return s.store.Update(upd)
}

func (s *service) GetNewsById(cmd *GetNewsByIdCommand) (*News, error) {
	return s.store.GetById(cmd.Id)
}

func (s *service) GetNewsByAuthorId(cmd *GetNewsByAuthorId) ([]News, error) {
	return s.store.GetByAuthorId(cmd.AuthorId)
}

func (s *service) UploadPhoto(cmd *UploadPhotoCommand) (*UploadPhotoResponse, error) {
	if cmd.ContentType == "" {
		return nil, ErrCannotDetectContentType
	}
	if !common.IsImage(cmd.ContentType) && !common.IsVideo(cmd.ContentType) {
		return nil, ErrFileShouldBeOnlyImageOrVideo
	}
	modelUpdate := &UpdateNews{
		Id: cmd.Id,
	}
	fileType := strings.Split(cmd.ContentType, "/")[1]
	fileResponse, err := s.s3Uploader.UploadFile(cmd.File.Bytes(), cmd.Id, fileType, cmd.ContentType)
	if err != nil {
		return nil, err
	}
	if common.IsImage(cmd.ContentType) {
		modelUpdate.PhotoUrl = &fileResponse.FileUrl
	} else {
		return nil, ErrInCheckForContentType
	}
	_, err = s.store.Update(modelUpdate)
	if err != nil {
		return nil, err
	}
	response := &UploadPhotoResponse{FileUrl: fileResponse.FileUrl}
	return response, nil
}
