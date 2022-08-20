package auth

type TokenStore interface {
	CreateToken(cmd *CreateTokenCommand) (*TokenDetails, error)
	GetToken(id string) (string, error)
	RemoveToken(id string) (int64, error)
}
