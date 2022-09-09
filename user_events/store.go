package user_events

type Store interface {
	Create(obj *UserEvent) (*UserEvent, error)
	ListByEventId(id string) ([]UserEvent, error)
	ListByUserId(id string) ([]UserEvent, error)
	List() ([]UserEvent, error)
	GetUserEventById(id string) (*UserEvent, error)
}
