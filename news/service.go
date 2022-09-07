package news

type Service interface {
	CreateNews(cmd *CreateNewsCommand) (*News, error)
	ListNews(cmd *ListNewsCommand) ([]News, error)
	UpdateNews(cmd *UpdateNewsCommand) (*News, error)
	GetNewsById(cmd *GetNewsByIdCommand) (*News, error)
	GetNewsByAuthorId(cmd *GetNewsByAuthorId) ([]News, error)
}
