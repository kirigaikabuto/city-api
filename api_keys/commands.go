package api_keys

type CreateApiKeyCommand struct {
}

func (cmd *CreateApiKeyCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(ApiKeyStore).Create()
}

type ListApiKeyCommand struct {
}

func (cmd *ListApiKeyCommand) Exec(svc interface{}) (interface{}, error) {
	return svc.(ApiKeyStore).List()
}
