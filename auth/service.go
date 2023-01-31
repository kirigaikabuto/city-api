package auth

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kirigaikabuto/city-api/common"
	"github.com/kirigaikabuto/city-api/mdw"
	sms_store "github.com/kirigaikabuto/city-api/sms-store"
	"github.com/kirigaikabuto/city-api/users"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"math/rand"
	"strings"
	"time"
)

type Service interface {
	Login(cmd *LoginCommand) (*mdw.Token, error)
	Register(cmd *RegisterCommand) (*users.User, error)
	GetProfile(cmd *GetMyProfileCommand) (*users.User, error)
	UpdateProfile(cmd *UpdateProfileCommand) (*users.User, error)
	UploadAvatar(cmd *UploadAvatarCommand) (*UploadAvatarResponse, error)
	VerifyCode(cmd *VerifyCodeCommand) error
	ResetPasswordRequest(cmd *ResetPasswordRequestCommand) error
	ResetPassword(cmd *ResetPasswordCommand) error
}

type service struct {
	userStore        users.UsersStore
	tokenStore       mdw.TokenStore
	s3               common.S3Uploader
	smsPostgresStore sms_store.Store
	smsTwilioStore   sms_store.Store
}

func NewService(u users.UsersStore, t mdw.TokenStore, s3 common.S3Uploader, smsPostgresStore sms_store.Store, smsTwilioStore sms_store.Store) Service {
	return &service{
		userStore:        u,
		tokenStore:       t,
		s3:               s3,
		smsPostgresStore: smsPostgresStore,
		smsTwilioStore:   smsTwilioStore,
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
	if !user.IsVerified {
		return nil, ErrUserIsNotVerified
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
	rand.Seed(time.Now().UnixNano())
	code, err := GenerateOTP(6)
	if err != nil {
		return nil, err
	}
	err = s.SendSmsCode(&SendSmsData{
		PhoneNumber: cmd.PhoneNumber,
		Title:       "Register template",
		Body:        code,
	})
	if err != nil {
		return nil, err
	}
	err = s.tokenStore.SaveCode(&mdw.SaveCodeCommand{
		Code:   code,
		UserId: user.Id,
		Time:   5 * time.Minute,
	})
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
	userUpdate.Id = cmd.Id
	if cmd.Username != nil && *cmd.Username != oldUser.Username {
		userUpdate.Username = cmd.Username
	}
	if cmd.Password != nil {
		hashedPassword, err := setdata_common.HashPassword(*cmd.Password)
		if err != nil {
			return nil, err
		}
		if oldUser.Password != hashedPassword {
			userUpdate.Password = &hashedPassword
		}
	}
	if cmd.FirstName != nil && *cmd.FirstName != oldUser.FirstName {
		userUpdate.FirstName = cmd.FirstName
	}
	if cmd.LastName != nil && *cmd.LastName != oldUser.LastName {
		userUpdate.LastName = cmd.LastName
	}
	if cmd.Email != nil && *cmd.Email != oldUser.Email {
		userUpdate.Email = cmd.Email
	}
	if cmd.PhoneNumber != nil && *cmd.PhoneNumber != oldUser.PhoneNumber {
		userUpdate.PhoneNumber = cmd.PhoneNumber
	}
	if cmd.Gender != nil && *cmd.Gender != oldUser.Gender {
		userUpdate.Gender = cmd.Gender
	}
	if cmd.Avatar != nil && *cmd.Avatar != oldUser.Avatar {
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

func (s *service) VerifyCode(cmd *VerifyCodeCommand) error {
	userId, err := s.tokenStore.GetUserIdByCode(cmd.Code)
	if err != nil && err == redis.Nil {
		return ErrInvalidCode
	} else if err != nil {
		return err
	}
	boolVal := true
	_, err = s.userStore.Update(&users.UserUpdate{IsVerified: &boolVal, Id: userId})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPasswordRequest(cmd *ResetPasswordRequestCommand) error {
	isPhone := false
	var user *users.User
	var err error
	if cmd.Email != "" {
		user = nil
	} else if cmd.PhoneNumber != "" {
		user, err = s.userStore.GetByPhoneNumber(cmd.PhoneNumber)
		if err != nil {
			return err
		}
		isPhone = true
	}
	code, err := GenerateOTP(6)
	if err != nil {
		return err
	}
	code = "111111"
	err = s.tokenStore.SaveCode(&mdw.SaveCodeCommand{
		Code:   code,
		UserId: user.Id,
		Time:   10 * time.Minute,
	})
	if err != nil {
		return err
	}
	if isPhone {
		err = s.SendSmsCode(&SendSmsData{
			PhoneNumber: cmd.PhoneNumber,
			Title:       "reset password",
			Body:        code,
		})
	}
	return nil
}

func GenerateOTP(length int) (string, error) {
	otpChars := "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}

func (s *service) SendSmsCode(data *SendSmsData) error {
	fmt.Println(data.Body)
	_, err := s.smsPostgresStore.Create(&sms_store.SmsCode{
		Title: data.Title,
		To:    data.PhoneNumber,
		Body:  data.Body,
	})
	if err != nil {
		return err
	}
	_, err = s.smsTwilioStore.Create(&sms_store.SmsCode{
		Title: data.Title,
		To:    data.PhoneNumber,
		Body:  data.Body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPassword(cmd *ResetPasswordCommand) error {
	userId, err := s.tokenStore.GetUserIdByCode(cmd.Code)
	if err != nil && err == redis.Nil {
		return ErrInvalidCode
	} else if err != nil {
		return err
	}
	user, err := s.userStore.Get(userId)
	if err != nil {
		return err
	}
	_, err = s.UpdateProfile(&UpdateProfileCommand{users.UserUpdate{Id: user.Id, Password: &cmd.NewPassword}})
	if err != nil {
		return err
	}
	return nil
}
