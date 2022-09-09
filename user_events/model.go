package user_events

type UserEvent struct {
	Id          string `json:"id"`
	EventId     string `json:"event_id"`
	UserId      string `json:"user_id"`
	CreatedDate string `json:"created_date"`
}
