package mdw

type TokenStore interface {
	CreateToken(cmd *CreateTokenCommand) (*TokenDetails, error)
	SaveCode(cmd *SaveCodeCommand) error
	GetUserIdByCode(code string) (string, error)
	GetToken(id string) (string, error)
	RemoveToken(id string) (int64, error)
	SaveApiToken(cmd *SaveApiTokenCommand) error
	GetApiToken(cmd *GetApiTokenCommand) (string, error)
}
