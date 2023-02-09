package sms_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kirigaikabuto/city-api/common"
	"io"
	"net/http"
)

type EmailStore struct {
	emailFrom    string
	ClientId     string
	ClientSecret string
	basicUrl     string
	html         string
}

func NewPulseEmailStore(config common.PulseEmailConfig) *EmailStore {
	return &EmailStore{
		emailFrom:    config.EmailFrom,
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		basicUrl:     config.BasicUrl,
		html:         "PHA+RXhhbXBsZSB0ZXh0PC9wPg==",
	}
}

func (t *EmailStore) SendEmail(req *SendEmailReq) error {
	fmt.Println("basic url", t.basicUrl)
	var to []UserObj
	to = append(to, UserObj{
		Name:  req.ToName,
		Email: req.ToEmail,
	})
	body := &SendEmailHttpReq{Email: EmailObj{
		Text:    req.Body,
		Subject: req.Subject,
		From: UserObj{
			Name:  "Чистый город",
			Email: t.emailFrom,
		},
		To: to,
	}}
	jsonBody, err := json.Marshal(body)
	fmt.Println(string(jsonBody))
	if err != nil {
		return err
	}
	finalBody := bytes.NewReader(jsonBody)
	request, err := http.NewRequest("POST", t.basicUrl+"/smtp/emails", finalBody)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+req.AccessToken)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	fmt.Println("client resp", res)
	return nil
}

func (t *EmailStore) GetToken() (*GetTokenResponse, error) {
	req := &GetTokenReq{
		ClientId:     t.ClientId,
		ClientSecret: t.ClientSecret,
		GrantType:    "client_credentials",
	}
	fmt.Println("basic url", t.basicUrl)
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	finalBody := bytes.NewReader(jsonBody)
	request, err := http.NewRequest("POST", t.basicUrl+"/oauth/access_token", finalBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	response := &GetTokenResponse{}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: response: %d\n", response)
	return response, nil
}
