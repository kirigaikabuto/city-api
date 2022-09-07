package news

import "bytes"

type CreateNewsCommand struct {
	News
}

func (cmd *CreateNewsCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).CreateNews(cmd)
}

type UpdateNewsCommand struct {
	UpdateNews
}

func (cmd *UpdateNewsCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UpdateNews(cmd)
}

type GetNewsByIdCommand struct {
	Id string `json:"id"`
}

func (cmd *GetNewsByIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetNewsById(cmd)
}

type ListNewsCommand struct {
}

func (cmd *ListNewsCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListNews(cmd)
}

type GetNewsByAuthorId struct {
	AuthorId string `json:"author_id"`
}

func (cmd *GetNewsByAuthorId) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).GetNewsByAuthorId(cmd)
}

type UploadPhotoCommand struct {
	UserId      string        `json:"-"`
	Id          string        `json:"id"`
	File        *bytes.Buffer `json:"file" form:"file"`
	ContentType string        `json:"-"`
}

func (cmd *UploadPhotoCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).UploadPhoto(cmd)
}

type UploadPhotoResponse struct {
	FileUrl string `json:"file_url"`
}
