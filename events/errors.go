package events

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateEventUnknown = com.NewMiddleError(errors.New("could not create event: unknown error"), 500, 51)
)
