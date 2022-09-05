package auth

type LoginCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cmd *LoginCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Login(cmd)
}

type RegisterCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cmd *RegisterCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(Service).Register(cmd)
}
