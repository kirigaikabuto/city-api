package file_storage

type Store interface {
	Create(model *FileStorage) (*FileStorage, error)
	ListByObjectType(objType string) ([]FileStorage, error)
	ListByObjectId(objId string) ([]FileStorage, error)
	List() ([]FileStorage, error)
}
