package news

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateNewsUnknown = com.NewMiddleError(errors.New("could not create news: unknown error"), 500, 51)
	ErrNoUserIdInToken   = com.NewMiddleError(errors.New("no user id in token"), 500, 52)
	ErrNewsNotFound      = com.NewMiddleError(errors.New("news not found"), 500, 53)
	ErrNothingToUpdate   = com.NewMiddleError(errors.New("nothing to update"), 500, 4)
)
