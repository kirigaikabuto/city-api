package user_events

type CreateUserEventCommand struct {
	UserEvent
}

func (cmd *CreateUserEventCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateUserEvent(cmd)
}

type ListByEventIdCommand struct {
	EventId string `json:"event_id"`
}

func (cmd *ListByEventIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListByEventId(cmd)
}

type ListByUserIdCommand struct {
	UserId string `json:"user_id"`
}

func (cmd *ListByUserIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListByUserId(cmd)
}

type ListUserEventsCommand struct {
}

func (cmd *ListUserEventsCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListUserEvents(cmd)
}

type GetUserEventByIdCommand struct {
	Id string `json:"id"`
}

func (cmd *GetUserEventByIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetUserEventById(cmd)
}
