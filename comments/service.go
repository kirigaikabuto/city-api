package comments

import (
	"github.com/kirigaikabuto/city-api/applications"
	"github.com/kirigaikabuto/city-api/events"
)

type Service interface {
	Create(cmd *CreateCommand) (*Comment, error)
	List(cmd *ListCommand) ([]Comment, error)
	ListByObjType(cmd *ListByObjTypeCommand) ([]Comment, error)
	ListByObjectId(cmd *ListByObjectIdCommand) ([]Comment, error)
}

type service struct {
	store            Store
	eventStore       events.Store
	applicationStore applications.Store
}

func NewService(s Store, eventStore events.Store, applicationStore applications.Store) Service {
	return &service{store: s, eventStore: eventStore, applicationStore: applicationStore}
}

func (s *service) Create(cmd *CreateCommand) (*Comment, error) {
	if !IsObjTypeExist(cmd.ObjType) {
		return nil, ErrObjTypeIncorrect
	}
	if cmd.ObjType == ApplicationObjType.ToString() {
		_, err := s.applicationStore.GetById(cmd.ObjId)
		if err != nil {
			return nil, err
		}
	} else if cmd.ObjType == EventObjType.ToString() {
		_, err := s.eventStore.GetById(cmd.ObjId)
		if err != nil {
			return nil, err
		}
	}
	return s.store.Create(&Comment{
		Message: cmd.Message,
		UserId:  cmd.UserId,
		ObjId:   cmd.ObjId,
		ObjType: ToObjType(cmd.ObjType),
	})
}

func (s *service) List(cmd *ListCommand) ([]Comment, error) {
	return s.store.List()
}

func (s *service) ListByObjType(cmd *ListByObjTypeCommand) ([]Comment, error) {
	if !IsObjTypeExist(cmd.ObjType) {
		return nil, ErrObjTypeIncorrect
	}
	return s.store.GetByObjType(ToObjType(cmd.ObjType))
}

func (s *service) ListByObjectId(cmd *ListByObjectIdCommand) ([]Comment, error) {
	comments, err := s.store.GetByObjId(cmd.ObjectId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
