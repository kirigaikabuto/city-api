package applications

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateApplicationUnknown     = com.NewMiddleError(errors.New("could not create applications: unknown error"), 500, 1)
	ErrApplicationTypeNotExist      = com.NewMiddleError(errors.New("applications type not exist"), 500, 2)
	ErrSearchPlaceNoAddressName     = com.NewMiddleError(errors.New("no address name"), 500, 3)
	ErrNothingToUpdate              = com.NewMiddleError(errors.New("nothing to update"), 500, 4)
	ErrApplicationNotFound          = com.NewMiddleError(errors.New("application not found"), 500, 5)
	ErrNoApplicationId              = com.NewMiddleError(errors.New("no application id in query"), 500, 6)
	ErrNoApplicationType            = com.NewMiddleError(errors.New("no application type in query"), 500, 7)
	ErrCannotDetectContentType      = com.NewMiddleError(errors.New("cannot detect content type"), 500, 8)
	ErrInCheckForContentType        = com.NewMiddleError(errors.New("check for content type incorrect"), 500, 9)
	ErrFileShouldBeOnlyImageOrVideo = com.NewMiddleError(errors.New("file should be image or video"), 500, 10)
	ErrApplicationStatusNotExist    = com.NewMiddleError(errors.New("applications status not exist"), 500, 11)
	ErrNoUserIdInToken              = com.NewMiddleError(errors.New("no user id in token"), 500, 12)
)
