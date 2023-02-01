package file_storage

import (
	"errors"
	com "github.com/kirigaikabuto/setdata-common"
)

var (
	ErrCreateUnknown = com.NewMiddleError(errors.New("could not file storage event: unknown error"), 500, 51)
)
