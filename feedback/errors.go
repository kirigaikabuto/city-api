package feedback

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateFeedbackUnknown = com.NewMiddleError(errors.New("could not create feedback: unknown error"), 500, 51)
	ErrNoUserIdInToken       = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrFeedbackNotFound      = com.NewMiddleError(errors.New("feedback not found"), 500, 53)
)
