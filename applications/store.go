package applications

type Store interface {
	Create(model *Application) (*Application, error)
	List() ([]Application, error)
	GetById(id string) (*Application, error)
	GetByProblemType(problemType ProblemType) ([]Application, error)
	Update(model *ApplicationUpdate) (*Application, error)
	ListApplicationsByUserId(userId string) ([]Application, error)
	RemoveApplication(id string) error
	ListByAddress(address string) ([]Application, error)
}
