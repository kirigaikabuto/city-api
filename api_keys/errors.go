package api_keys

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateApiKeyUnknown = com.NewMiddleError(errors.New("could not create api key: unknown error"), 500, 101)
	ErrApiKeyNotFound      = com.NewMiddleError(errors.New("api key not found"), 500, 102)
)
