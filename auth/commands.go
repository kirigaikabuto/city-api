package auth

import (
	"bytes"
	"github.com/kirigaikabuto/city-api/users"
)

type LoginCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cmd *LoginCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Login(cmd)
}

type RegisterCommand struct {
	users.User
}

func (cmd *RegisterCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Register(cmd)
}

type GetMyProfileCommand struct {
	UserId string `json:"-"`
}

func (cmd *GetMyProfileCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetProfile(cmd)
}

type UpdateProfileCommand struct {
	users.UserUpdate
}

func (cmd *UpdateProfileCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UpdateProfile(cmd)
}

type UploadAvatarCommand struct {
	UserId      string        `json:"user_id"`
	File        *bytes.Buffer `json:"file" form:"file"`
	ContentType string        `json:"-"`
}

func (cmd *UploadAvatarCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UploadAvatar(cmd)
}

type UploadAvatarResponse struct {
	FileUrl string `json:"file_url"`
}

type VerifyCodeCommand struct {
	Code string `json:"code"`
}

func (cmd *VerifyCodeCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(Service).VerifyCode(cmd)
}

type ResetPasswordRequestCommand struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (cmd *ResetPasswordRequestCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(Service).ResetPasswordRequest(cmd)
}

type SendEmailData struct {
	ToEmail     string
	ToName      string
	Body        string
	Subject     string
	AccessToken string
}

type ResetPasswordCommand struct {
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func (cmd *ResetPasswordCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(Service).ResetPassword(cmd)
}

type RemoveAccountCommand struct {
	UserId string `json:"user_id"`
}

func (cmd *RemoveAccountCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(Service).RemoveAccount(cmd)
}
