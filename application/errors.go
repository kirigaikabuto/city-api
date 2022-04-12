package application

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateApplicationUnknown = com.NewMiddleError(errors.New("could not create application: unknown error"), 500, 1)
	ErrApplicationTypeNotExist  = com.NewMiddleError(errors.New("application type not exist"), 500, 2)
	ErrSearchPlaceNoAddressName = com.NewMiddleError(errors.New("no address name"), 500, 3)
)
