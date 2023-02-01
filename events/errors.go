package events

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateEventUnknown           = com.NewMiddleError(errors.New("could not create event: unknown error"), 500, 51)
	ErrNoUserIdInToken              = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrNothingToUpdate              = com.NewMiddleError(errors.New("nothing to update"), 500, 4)
	ErrEventNotFound                = com.NewMiddleError(errors.New("event not found"), 500, 5)
	ErrNoEventId                    = com.NewMiddleError(errors.New("no event id"), 500, 7)
	ErrCannotDetectContentType      = com.NewMiddleError(errors.New("cannot detect content type"), 500, 8)
	ErrInCheckForContentType        = com.NewMiddleError(errors.New("check for content type incorrect"), 500, 9)
	ErrFileShouldBeOnlyImageOrVideo = com.NewMiddleError(errors.New("file should be image or video"), 500, 10)
)
