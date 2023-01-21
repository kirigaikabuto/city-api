package sms_store

type Store interface {
	Create(obj *SmsCode) (*SmsCode, error)
	GetById(id string) (*SmsCode, error)
}
