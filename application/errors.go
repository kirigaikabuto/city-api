package application

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateApplicationUnknown = com.NewMiddleError(errors.New("could not create application: unknown error"), 500, 1)
)
