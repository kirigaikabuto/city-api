package feedback

type CreateFeedbackCommand struct {
	Feedback
}

func (cmd *CreateFeedbackCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateFeedback(cmd)
}

type ListFeedbackCommand struct {
}

func (cmd *ListFeedbackCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListFeedback(cmd)
}
