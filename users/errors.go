package users

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrUserNotFound                 = com.NewMiddleError(errors.New("user not found"), 404, 201)
	ErrUserIdNotProvided            = com.NewMiddleError(errors.New("user id is not provided"), 400, 202)
	ErrUserPasswordNotCorrect       = com.NewMiddleError(errors.New("user password not correct"), 500, 203)
	ErrCreateUserUnknown            = com.NewMiddleError(errors.New("could not create user: unknown error"), 500, 204)
	ErrNothingToUpdate              = com.NewMiddleError(errors.New("nothing to update"), 400, 205)
	ErrUserWithUsernameAlreadyExist = com.NewMiddleError(errors.New("user with that username already exist"), 400, 206)
)
