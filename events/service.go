package events

type Service interface {
	Create(cmd *CreateEventCommand) (*Event, error)
	List(cmd *ListEventCommand) ([]Event, error)
}

type service struct {
	eventStore Store
}

func NewService(e Store) Service {
	return &service{eventStore: e}
}

func (s *service) Create(cmd *CreateEventCommand) (*Event, error) {
	return s.eventStore.Create(cmd.Event)
}

func (s *service) List(cmd *ListEventCommand) ([]Event, error) {
	return s.eventStore.List()
}
