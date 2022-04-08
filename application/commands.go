package application

type CreateApplicationCommand struct {
	Address     string `json:"address"`
	AppType     string `json:"app_type"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Patronymic  string `json:"patronymic"`
	PhoneNumber string `json:"phone_number"`
	PhotoUrl    string `json:"photo_url"`
	VideoUrl    string `json:"video_url"`
	UserId      string `json:"-"`
}

func (cmd *CreateApplicationCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateApplication(cmd)
}
