package user_events

type UserEvent struct {
	Id          string  `json:"id"`
	ObjId       string  `json:"obj_id"`
	ObjType     ObjType `json:"obj_type"`
	UserId      string  `json:"user_id"`
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
