package user_events

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateUserEventUnknown = com.NewMiddleError(errors.New("could not create user event: unknown error"), 500, 51)
	ErrNoUserIdInToken        = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrUserEventNotFound      = com.NewMiddleError(errors.New("user event not found"), 500, 53)
	ErrNoObjType              = com.NewMiddleError(errors.New("obj type can be only application/event"), 400, 54)
)
