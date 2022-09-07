package news

type Store interface {
	Create(obj *News) (*News, error)
	Update(obj *UpdateNews) (*News, error)
	List() ([]News, error)
	GetById(id string) (*News, error)
	GetByAuthorId(authorId string) ([]News, error)
}
