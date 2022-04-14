package users

type CreateUserCommand struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (cmd *CreateUserCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(UserService).CreateUser(cmd)
}

type UpdateUserCommand struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (cmd *UpdateUserCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(UserService).UpdateUser(cmd)
}

type DeleteUserCommand struct {
	Id string `json:"id"`
}

func (cmd *DeleteUserCommand) Exec(svc interface{}) (interface{}, error) {
	return nil, svc.(UserService).DeleteUser(cmd)
}

type GetUserCommand struct {
	Id string `json:"id"`
}

func (cmd *GetUserCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(UserService).GetUser(cmd)
}

type ListUserCommand struct {
}

func (cmd *ListUserCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(UserService).ListUser(cmd)
}

type GetUserByUsernameAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cmd *GetUserByUsernameAndPassword) Exec(svc interface{}) (interface{}, error) {
	return svc.(UserService).GetUserByUsernameAndPassword(cmd)
}
