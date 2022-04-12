package events

type Store interface {
	Create(model *Event) (*Event, error)
	List() ([]Event, error)
}
