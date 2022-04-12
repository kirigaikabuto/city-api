package events

type Event struct {
	Id            string  `json:"id"`
	Address       string  `json:"address"`
	Description   string  `json:"description"`
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	OrganizerInfo string  `json:"organizer_info"`
	DocumentUrl   string  `json:"document_url"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	CreatedDate   string  `json:"created_date"`
}
