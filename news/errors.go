package news

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateNewsUnknown            = com.NewMiddleError(errors.New("could not create news: unknown error"), 500, 51)
	ErrNoUserIdInToken              = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrNewsNotFound                 = com.NewMiddleError(errors.New("news not found"), 500, 53)
	ErrNothingToUpdate              = com.NewMiddleError(errors.New("nothing to update"), 500, 4)
	ErrCannotDetectContentType      = com.NewMiddleError(errors.New("cannot detect content type"), 500, 8)
	ErrInCheckForContentType        = com.NewMiddleError(errors.New("check for content type incorrect"), 500, 9)
	ErrFileShouldBeOnlyImageOrVideo = com.NewMiddleError(errors.New("file should be image or video"), 500, 10)
)
