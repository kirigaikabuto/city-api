package auth

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrPleaseFillUsernamePassword = com.NewMiddleError(errors.New("please fill username and password"), 400, 200)
)
