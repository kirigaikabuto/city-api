package feedback

type Feedback struct {
	Id          string `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	CreatedDate string `json:"created_date"`
}
