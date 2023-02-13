package user_events

import (
	"github.com/kirigaikabuto/city-api/applications"
	"github.com/kirigaikabuto/city-api/events"
	"github.com/kirigaikabuto/city-api/users"
)

type Service interface {
	CreateUserEvent(cmd *CreateUserEventCommand) (*UserEvent, error)
	ListByObjectId(cmd *ListByEventIdCommand) ([]UserEvent, error)
	ListByUserId(cmd *ListByUserIdCommand) ([]UserEvent, error)
	ListUserEvents(cmd *ListUserEventsCommand) ([]UserEvent, error)
	GetUserEventById(cmd *GetUserEventByIdCommand) (*UserEvent, error)
}

type service struct {
	store       Store
	usersStore  users.UsersStore
	eventsStore events.Store
	appStore    applications.Store
}

func NewService(s Store, u users.UsersStore, e events.Store, a applications.Store) Service {
	return &service{
		store:       s,
		usersStore:  u,
		eventsStore: e,
		appStore:    a,
	}
}

func (s *service) CreateUserEvent(cmd *CreateUserEventCommand) (*UserEvent, error) {
	_, err := s.usersStore.Get(cmd.UserId)
	if err != nil {
		return nil, err
	}
	if cmd.ObjType == EventObjType {
		_, err = s.eventsStore.GetById(cmd.ObjId)
		if err != nil {
			return nil, err
		}
	} else if cmd.ObjType == ApplicationObjType {
		_, err := s.appStore.GetById(cmd.ObjId)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrNoObjType
	}
	return s.store.Create(&cmd.UserEvent)
}

func (s *service) ListByObjectId(cmd *ListByEventIdCommand) ([]UserEvent, error) {
	return s.store.ListByObjectId(cmd.EventId)
}

func (s *service) ListByUserId(cmd *ListByUserIdCommand) ([]UserEvent, error) {
	return s.store.ListByUserId(cmd.UserId)
}

func (s *service) ListUserEvents(cmd *ListUserEventsCommand) ([]UserEvent, error) {
	return s.store.List()
}

func (s *service) GetUserEventById(cmd *GetUserEventByIdCommand) (*UserEvent, error) {
	return s.store.GetUserEventById(cmd.Id)
}
