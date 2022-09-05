package auth

import (
	"github.com/kirigaikabuto/city-api/common"
	"github.com/kirigaikabuto/city-api/mdw"
	"github.com/kirigaikabuto/city-api/users"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"strings"
)

type Service interface {
	Login(cmd *LoginCommand) (*mdw.Token, error)
	Register(cmd *RegisterCommand) (*users.User, error)
	GetProfile(cmd *GetMyProfileCommand) (*users.User, error)
	UpdateProfile(cmd *UpdateProfileCommand) (*users.User, error)
	UploadAvatar(cmd *UploadAvatarCommand) (*UploadAvatarResponse, error)
}

type service struct {
	userStore  users.UsersStore
	tokenStore mdw.TokenStore
	s3         common.S3Uploader
}

func NewService(u users.UsersStore, t mdw.TokenStore) Service {
	return &service{
		userStore:  u,
		tokenStore: t,
	}
}

func (s *service) Login(cmd *LoginCommand) (*mdw.Token, error) {
	if cmd.Username == "" || cmd.Password == "" {
		return nil, ErrPleaseFillUsernamePassword
	}
	user, err := s.userStore.GetByUsernameAndPassword(cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}
	tkn, err := s.tokenStore.CreateToken(&mdw.CreateTokenCommand{
		UserId:   user.Id,
		UserType: user.AccessType.ToString(),
	})
	return &mdw.Token{
		AccessToken:  tkn.AccessToken,
		RefreshToken: tkn.RefreshToken,
	}, err
}

func (s *service) Register(cmd *RegisterCommand) (*users.User, error) {
	if cmd.Username == "" || cmd.Password == "" {
		return nil, ErrPleaseFillUsernamePassword
	}
	if !users.IsGenderExist(cmd.Gender.ToString()) {
		return nil, ErrNoGenderType
	}
	cmd.AccessType = users.AccessTypeUser
	user, err := s.userStore.Create(&cmd.User)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetProfile(cmd *GetMyProfileCommand) (*users.User, error) {
	return s.userStore.Get(cmd.UserId)
}

func (s *service) UpdateProfile(cmd *UpdateProfileCommand) (*users.User, error) {
	userUpdate := &users.UserUpdate{}
	oldUser, err := s.userStore.Get(cmd.Id)
	if err != nil {
		return nil, err
	}
	if *cmd.Username != "" && *cmd.Username != oldUser.Username {
		userUpdate.Username = cmd.Username
	}
	if *cmd.Password != "" {
		hashedPassword, err := setdata_common.HashPassword(*cmd.Password)
		if err != nil {
			return nil, err
		}
		if oldUser.Password != hashedPassword {
			userUpdate.Password = &hashedPassword
		}
	}
	if *cmd.FirstName != "" && *cmd.FirstName != oldUser.FirstName {
		userUpdate.FirstName = cmd.FirstName
	}
	if *cmd.LastName != "" && *cmd.LastName != oldUser.LastName {
		userUpdate.LastName = cmd.LastName
	}
	if *cmd.Email != "" && *cmd.Email != oldUser.Email {
		userUpdate.Email = cmd.Email
	}
	if *cmd.PhoneNumber != "" && *cmd.PhoneNumber != oldUser.PhoneNumber {
		userUpdate.PhoneNumber = cmd.PhoneNumber
	}
	if *cmd.Gender != "" && *cmd.Gender != oldUser.Gender {
		userUpdate.Gender = cmd.Gender
	}
	if *cmd.Avatar != "" && *cmd.Avatar != oldUser.Avatar {
		userUpdate.Avatar = cmd.Avatar
	}
	return s.userStore.Update(userUpdate)
}

func (s *service) UploadAvatar(cmd *UploadAvatarCommand) (*UploadAvatarResponse, error) {
	if cmd.ContentType == "" {
		return nil, ErrCannotDetectContentType
	}
	if !common.IsImage(cmd.ContentType) && !common.IsVideo(cmd.ContentType) {
		return nil, ErrFileShouldBeOnlyImageOrVideo
	}
	modelUpdate := &UpdateProfileCommand{}
	modelUpdate.Id = cmd.UserId
	fileType := strings.Split(cmd.ContentType, "/")[1]
	fileResponse, err := s.s3.UploadFile(cmd.File.Bytes(), cmd.UserId, fileType, cmd.ContentType)
	if err != nil {
		return nil, err
	}
	if common.IsImage(cmd.ContentType) {
		modelUpdate.Avatar = &fileResponse.FileUrl
	} else {
		return nil, ErrInCheckForContentType
	}
	_, err = s.UpdateProfile(modelUpdate)
	if err != nil {
		return nil, err
	}
	response := &UploadAvatarResponse{FileUrl: fileResponse.FileUrl}
	return response, nil
}
