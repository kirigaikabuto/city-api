package users

import (
	"github.com/google/uuid"
	setdata_common "github.com/kirigaikabuto/setdata-common"
)

type UserService interface {
	CreateUser(cmd *CreateUserCommand) (*User, error)
	UpdateUser(cmd *UpdateUserCommand) (*User, error)
	DeleteUser(cmd *DeleteUserCommand) error
	GetUser(cmd *GetUserCommand) (*User, error)
	ListUser(cmd *ListUserCommand) ([]User, error)
	GetUserByUsernameAndPassword(cmd *GetUserByUsernameAndPassword) (*User, error)
}

type userService struct {
	userStore UsersStore
}

func NewUserService(uStore UsersStore) UserService {
	return &userService{userStore: uStore}
}

func (u *userService) CreateUser(cmd *CreateUserCommand) (*User, error) {
	user := &User{
		Id:       uuid.New().String(),
		Username: cmd.Username,
		Password: cmd.Password,
	}
	return u.userStore.Create(user)
}

func (u *userService) UpdateUser(cmd *UpdateUserCommand) (*User, error) {
	if cmd.Id == "" {
		return nil, ErrUserIdNotProvided
	}
	oldUser, err := u.GetUser(&GetUserCommand{
		Id: cmd.Id,
	})
	if err != nil {
		return nil, err
	}
	userUpdate := &UserUpdate{Id: cmd.Id}
	if cmd.Password != "" {
		hashedPassword, err := setdata_common.HashPassword(cmd.Password)
		if err != nil {
			return nil, err
		}
		if oldUser.Password != hashedPassword {
			userUpdate.Password = &hashedPassword
		}
	}
	if cmd.Username != "" && cmd.Username != oldUser.Username {
		userUpdate.Username = &cmd.Username
	}
	return u.userStore.Update(userUpdate)
}

func (u *userService) DeleteUser(cmd *DeleteUserCommand) error {
	return u.userStore.Delete(cmd.Id)
}

func (u *userService) GetUser(cmd *GetUserCommand) (*User, error) {
	return u.userStore.Get(cmd.Id)
}

func (u *userService) ListUser(cmd *ListUserCommand) ([]User, error) {
	return u.userStore.List()
}

func (u *userService) GetUserByUsernameAndPassword(cmd *GetUserByUsernameAndPassword) (*User, error) {
	return u.userStore.GetByUsernameAndPassword(cmd.Username, cmd.Password)
}
