package api_keys

type ApiKeyStore interface {
	Create() (*ApiKey, error)
	GetByKey(key string) (*ApiKey, error)
	Delete(id string) error
	List() ([]ApiKey, error)
}
