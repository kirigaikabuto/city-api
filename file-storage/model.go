package file_storage

type FileStorage struct {
	Id          string  `json:"id"`
	ObjectId    string  `json:"object_id"`
	ObjectType  ObjType `json:"object_type"`
	FileUrl     string  `json:"file_url"`
	CreatedDate string  `json:"created_date"`
}

type ObjType string

var (
	ApplicationObjType ObjType = "application"
	EventObjType       ObjType = "event"
)

var (
	objTypeToString = map[ObjType]string{
		ApplicationObjType: "application",
		EventObjType:       "event",
	}
	stringToObjType = map[string]ObjType{
		"application": ApplicationObjType,
		"event":       EventObjType,
	}
)

func (c ObjType) ToString() string {
	return objTypeToString[c]
}

func ToObjType(s string) ObjType {
	return stringToObjType[s]
}

func IsObjTypeExist(s string) bool {
	objTypes := []string{"application", "event"}
	for _, v := range objTypes {
		if v == s {
			return true
		}
	}
	return false
}
