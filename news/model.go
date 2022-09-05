package news

type News struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PhotoUrl    string `json:"photo_url"`
}
