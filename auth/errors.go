package auth

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrNoApiKeyHeaderValue = com.NewMiddleError(errors.New("no api key header value"), 500, 151)
	ErrIncorrectApiKey     = com.NewMiddleError(errors.New("incorrect api key"), 500, 152)
)
