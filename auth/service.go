package auth

import (
	"github.com/kirigaikabuto/city-api/mdw"
	"github.com/kirigaikabuto/city-api/users"
)

type Service interface {
	Login(cmd *LoginCommand) (*mdw.Token, error)
	Register(cmd *RegisterCommand) (*users.User, error)
}

type service struct {
	userStore  users.UsersStore
	tokenStore mdw.TokenStore
}

func NewService(u users.UsersStore, t mdw.TokenStore) Service {
	return &service{
		userStore:  u,
		tokenStore: t,
	}
}

func (s *service) Login(cmd *LoginCommand) (*mdw.Token, error) {
	if cmd.Username == "" || cmd.Password == "" {
		return nil, ErrPleaseFillUsernamePassword
	}
	user, err := s.userStore.GetByUsernameAndPassword(cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}
	tkn, err := s.tokenStore.CreateToken(&mdw.CreateTokenCommand{
		UserId:   user.Username,
		UserType: user.AccessType.ToString(),
	})
	return &mdw.Token{
		AccessToken:  tkn.AccessToken,
		RefreshToken: tkn.RefreshToken,
	}, err
}

func (s *service) Register(cmd *RegisterCommand) (*users.User, error) {
	if cmd.Username == "" || cmd.Password == "" {
		return nil, ErrPleaseFillUsernamePassword
	}
	user, err := s.userStore.Create(&users.User{
		Username:   cmd.Username,
		Password:   cmd.Password,
		AccessType: users.AccessTypeUser,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
