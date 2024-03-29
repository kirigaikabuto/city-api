package mdw

import "time"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
	UserType   string
}

type CreateTokenCommand struct {
	UserId   string
	UserType string
}

type SaveCodeCommand struct {
	Code   string
	UserId string
	Time   time.Duration
}

type SaveApiTokenCommand struct {
	Value string
	Key   string
	Time  time.Duration
}

type GetApiTokenCommand struct {
	Key string
}
