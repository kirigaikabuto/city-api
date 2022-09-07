package news

type News struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	SmallDescription string `json:"small_description"`
	Description      string `json:"description"`
	PhotoUrl         string `json:"photo_url"`
	AuthorId         string `json:"author_id"`
	CreatedDate      string `json:"created_date"`
}

type UpdateNews struct {
	Id               string  `json:"id"`
	Title            *string `json:"title"`
	SmallDescription *string `json:"small_description"`
	Description      *string `json:"description"`
	PhotoUrl         *string `json:"photo_url"`
	AuthorId         *string `json:"author_id"`
	CreatedDate      *string `json:"created_date"`
}
