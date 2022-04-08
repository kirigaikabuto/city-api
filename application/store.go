package application

type Store interface {
	Create(model *Application) (*Application, error)
	List() ([]Application, error)
	GetById(id string) (*Application, error)
	GetByProblemType(problemType ProblemType) ([]Application, error)
}
