package events

type Store interface {
	Create(model *Event) (*Event, error)
	List() ([]Event, error)
	ListByUserId(userId string) ([]Event, error)
	Update(model *EventUpdate) (*Event, error)
	GetById(id string) (*Event, error)
}
