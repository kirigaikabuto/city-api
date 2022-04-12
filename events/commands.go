package events

type CreateEventCommand struct {
	*Event
}

func (cmd *CreateEventCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Create(cmd)
}

type ListEventCommand struct {
}

func (cmd *ListEventCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).List(cmd)
}
