package application

type Service interface {
	CreateApplication(cmd *CreateApplicationCommand) (*Application, error)
}

type service struct {
	appStore Store
}

func NewApplicationService(appStore Store) Service {
	return &service{appStore: appStore}
}

func (s *service) CreateApplication(cmd *CreateApplicationCommand) (*Application, error) {
	return nil, nil
}
