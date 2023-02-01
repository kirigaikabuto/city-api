package comments

import (
	"github.com/kirigaikabuto/city-api/applications"
	"github.com/kirigaikabuto/city-api/events"
	"github.com/kirigaikabuto/city-api/users"
)

type Service interface {
	Create(cmd *CreateCommand) (*Comment, error)
	List(cmd *ListCommand) ([]Comment, error)
	ListByObjType(cmd *ListByObjTypeCommand) ([]Comment, error)
	ListByObjectId(cmd *ListByObjectIdCommand) ([]ListByObjectIdResponse, error)
}

type service struct {
	store            Store
	eventStore       events.Store
	applicationStore applications.Store
	usersStore       users.UsersStore
}

func NewService(s Store, eventStore events.Store, applicationStore applications.Store, usersStore users.UsersStore) Service {
	return &service{store: s, eventStore: eventStore, applicationStore: applicationStore, usersStore: usersStore}
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

func (s *service) ListByObjectId(cmd *ListByObjectIdCommand) ([]ListByObjectIdResponse, error) {
	comments, err := s.store.GetByObjId(cmd.ObjectId)
	if err != nil {
		return nil, err
	}
	result := []ListByObjectIdResponse{}
	for _, v := range comments {
		temp := ListByObjectIdResponse{}
		temp.Comment = v
		user, err := s.usersStore.Get(v.UserId)
		if err != nil {
			return nil, err
		}
		temp.UserName = user.Username
		temp.UserPhoto = user.Avatar
		result = append(result, temp)
	}
	return result, nil
}
