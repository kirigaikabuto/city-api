package sms_store

import (
	"github.com/kirigaikabuto/city-api/common"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioStore struct {
	client *twilio.RestClient
}

func NewTwilioStore(config common.TwilioConfig) Store {
	client := twilio.NewRestClient()
	return &twilioStore{client: client}
}

func (t *twilioStore) Create(obj *SmsCode) (*SmsCode, error) {
	params := &api.CreateMessageParams{}
	params.SetBody(obj.Body)
	params.SetFrom(obj.From)
	params.SetTo(obj.To)
	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (t *twilioStore) GetById(id string) (*SmsCode, error) {
	return nil, nil
}
