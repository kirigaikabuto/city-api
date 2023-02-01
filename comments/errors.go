package comments

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateCommentUnknown = com.NewMiddleError(errors.New("could not create comments: unknown error"), 500, 51)
	ErrNoUserIdInToken      = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrCommentNotFound      = com.NewMiddleError(errors.New("comments not found"), 500, 53)
	ErrObjTypeIncorrect     = com.NewMiddleError(errors.New("no object type"), 500, 54)
	ErrNoObjectID           = com.NewMiddleError(errors.New("no object id"), 500, 55)
)
