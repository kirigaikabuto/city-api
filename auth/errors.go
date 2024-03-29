package auth

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrPleaseFillUsernamePassword   = com.NewMiddleError(errors.New("please fill username and password"), 400, 200)
	ErrNoUserIdInToken              = com.NewMiddleError(errors.New("no user id in token"), 500, 201)
	ErrNoGenderType                 = com.NewMiddleError(errors.New("no gender type"), 500, 202)
	ErrCannotDetectContentType      = com.NewMiddleError(errors.New("cannot detect content type"), 500, 203)
	ErrInCheckForContentType        = com.NewMiddleError(errors.New("check for content type incorrect"), 500, 204)
	ErrFileShouldBeOnlyImageOrVideo = com.NewMiddleError(errors.New("file should be image or video"), 500, 205)
	ErrNoCodeInQuery                = com.NewMiddleError(errors.New("no code in query"), 400, 206)
	ErrInvalidCode                  = com.NewMiddleError(errors.New("invalid code"), 400, 207)
	ErrUserIsNotVerified            = com.NewMiddleError(errors.New("user is not verified"), 400, 208)
	ErrNoPhoneNumberOrEmail         = com.NewMiddleError(errors.New("no phone number or email"), 400, 209)
	ErrNoPasswordInRequest          = com.NewMiddleError(errors.New("no password in request"), 400, 210)
	ErrAdminTypeDelete              = com.NewMiddleError(errors.New("access type admin can not be deleted"), 400, 211)
	ErrNotValidPhoneNumber          = com.NewMiddleError(errors.New("not valid phone number"), 500, 212)
	ErrUserWithUsernameAlreadyExist = com.NewMiddleError(errors.New("user with that username already exist"), 400, 213)
	ErrUserWithEmailAlreadyExist    = com.NewMiddleError(errors.New("user with that email already exist"), 400, 214)
	ErrPleaseFillEmail              = com.NewMiddleError(errors.New("please fill email"), 400, 215)
)
