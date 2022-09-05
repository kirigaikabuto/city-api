package comments

type Service interface {
	Create(cmd *CreateCommand) (*Comment, error)
	List(cmd *ListCommand) ([]Comment, error)
	ListByObjType(cmd *ListByObjTypeCommand) ([]Comment, error)
}

type service struct {
	store Store
}

func NewService(s Store) Service {
	return &service{store: s}
}

func (s *service) Create(cmd *CreateCommand) (*Comment, error) {
	if !IsObjTypeExist(cmd.ObjType) {
		return nil, ErrObjTypeIncorrect
	}
	return s.store.Create(&Comment{
		Message: cmd.Message,
		UserId:  cmd.UserId,
		ObjId:   cmd.ObjId,
		ObjType: ToObjType(cmd.ObjType),
	})
}

func (s *service) List(cmd *ListCommand) ([]Comment, error) {
	return s.store.List()
}

func (s *service) ListByObjType(cmd *ListByObjTypeCommand) ([]Comment, error) {
	if !IsObjTypeExist(cmd.ObjType) {
		return nil, ErrObjTypeIncorrect
	}
	return s.store.GetByObjType(ToObjType(cmd.ObjType))
}
