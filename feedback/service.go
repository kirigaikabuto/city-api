package feedback

type Service interface {
	CreateFeedback(cmd *CreateFeedbackCommand) (*Feedback, error)
	ListFeedback(cmd *ListFeedbackCommand) ([]Feedback, error)
}

type service struct {
	store Store
}

func NewService(s Store) Service {
	return &service{store: s}
}

func (s *service) CreateFeedback(cmd *CreateFeedbackCommand) (*Feedback, error) {
	return s.store.Create(&Feedback{
		Message:     cmd.Message,
		FullName:    cmd.FullName,
		PhoneNumber: cmd.PhoneNumber,
	})
}
func (s *service) ListFeedback(cmd *ListFeedbackCommand) ([]Feedback, error) {
	return s.store.List()
}
