package comments

type Store interface {
	Create(obj *Comment) (*Comment, error)
	GetById(id string) (*Comment, error)
	List() ([]Comment, error)
	GetByObjType(objType ObjType) ([]Comment, error)
}
