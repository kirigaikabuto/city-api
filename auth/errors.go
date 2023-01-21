package auth

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrPleaseFillUsernamePassword   = com.NewMiddleError(errors.New("please fill username and password"), 400, 200)
	ErrNoUserIdInToken              = com.NewMiddleError(errors.New("no user id in token"), 500, 201)
	ErrNoGenderType                 = com.NewMiddleError(errors.New("no gender type"), 500, 202)
	ErrCannotDetectContentType      = com.NewMiddleError(errors.New("cannot detect content type"), 500, 8)
	ErrInCheckForContentType        = com.NewMiddleError(errors.New("check for content type incorrect"), 500, 9)
	ErrFileShouldBeOnlyImageOrVideo = com.NewMiddleError(errors.New("file should be image or video"), 500, 10)
	ErrNoCodeInQuery                = com.NewMiddleError(errors.New("no code in query"), 400, 11)
	ErrInvalidCode                  = com.NewMiddleError(errors.New("invalid code"), 400, 12)
	ErrUserIsNotVerified            = com.NewMiddleError(errors.New("user is not verified"), 400, 13)
)
