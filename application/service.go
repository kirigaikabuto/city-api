package application

type Service interface {
	CreateApplication(cmd *CreateApplicationCommand) (*Application, error)
	ListApplications(cmd *ListApplicationsCommand) ([]Application, error)
}

type service struct {
	appStore Store
}

func NewApplicationService(appStore Store) Service {
	return &service{appStore: appStore}
}

func (s *service) CreateApplication(cmd *CreateApplicationCommand) (*Application, error) {
	var appType ProblemType
	if IsProblemTypeExist(cmd.AppType) {
		appType = ToProblemType(cmd.AppType)
	} else {
		return nil, ErrApplicationTypeNotExist
	}
	app := &Application{
		AppType:     appType,
		FirstName:   cmd.FirstName,
		LastName:    cmd.LastName,
		Patronymic:  cmd.Patronymic,
		PhoneNumber: cmd.PhoneNumber,
		Address:     cmd.Address,
		Longitude:   cmd.Longitude,
		Latitude:    cmd.Latitude,
	}
	return s.appStore.Create(app)
}

func (s *service) ListApplications(cmd *ListApplicationsCommand) ([]Application, error) {

}
