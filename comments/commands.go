package comments

type CreateCommand struct {
	Message string `json:"message"`
	UserId  string `json:"user_id"`
	ObjId   string `json:"obj_id"`
	ObjType string `json:"obj_type"`
}

func (cmd *CreateCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Create(cmd)
}

type ListCommand struct {
}

func (cmd *ListCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).List(cmd)
}

type ListByObjTypeCommand struct {
	ObjType string `json:"obj_type"`
}

func (cmd *ListByObjTypeCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListByObjType(cmd)
}

type ListByObjectIdCommand struct {
	ObjectId string `json:"object_id"`
}

func (cmd *ListByObjectIdCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).ListByObjectId(cmd)
}
