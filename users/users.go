package users

type User struct {
	Id         string     `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	AccessType AccessType `json:"access_type"`
}

type UserUpdate struct {
	Id         string      `json:"id"`
	Username   *string     `json:"username"`
	Password   *string     `json:"password"`
	AccessType *AccessType `json:"access_type"`
}

type AccessType string

var (
	AccessTypeUser  AccessType = "user"
	AccessTypeAdmin AccessType = "admin"
)

var (
	accessTypeToString = map[AccessType]string{
		AccessTypeUser:  "user",
		AccessTypeAdmin: "admin",
	}
	stringToAccessType = map[string]AccessType{
		"user":  AccessTypeUser,
		"admin": AccessTypeAdmin,
	}
)

func (c AccessType) ToString() string {
	return accessTypeToString[c]
}

func ToAccessType(s string) AccessType {
	return stringToAccessType[s]
}

func IsAccessTypeExist(s string) bool {
	accessTypes := []string{"user", "admin"}
	for _, v := range accessTypes {
		if v == s {
			return true
		}
	}
	return false
}
