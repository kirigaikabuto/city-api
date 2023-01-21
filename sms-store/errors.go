package sms_store

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateSmsUnknown = com.NewMiddleError(errors.New("could not create sms: unknown error"), 500, 51)
	ErrSmsNotFound      = com.NewMiddleError(errors.New("sms not found"), 500, 52)
)
