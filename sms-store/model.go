package sms_store

type SmsCode struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type GetTokenReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type UserObj struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailObj struct {
	Html    string    `json:"html"`
	Text    string    `json:"text"`
	Subject string    `json:"subject"`
	From    UserObj   `json:"from"`
	To      []UserObj `json:"to"`
}
type SendEmailHttpReq struct {
	Email EmailObj `json:"email"`
}

type SendEmailReq struct {
	ToEmail     string
	ToName      string
	Body        string
	Subject     string
	AccessToken string
}

type GetTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type CreateMessageResponse struct {
	Sid    string `json:"sid"`
	Status string `json:"status"`
}
