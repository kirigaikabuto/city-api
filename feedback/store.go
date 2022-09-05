package feedback

type Store interface {
	Create(obj *Feedback) (*Feedback, error)
	GetById(id string) (*Feedback, error)
	List() ([]Feedback, error)
}
