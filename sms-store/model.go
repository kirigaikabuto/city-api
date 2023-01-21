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

type CreateMessageResponse struct {
	Sid    string `json:"sid"`
	Status string `json:"status"`
}
