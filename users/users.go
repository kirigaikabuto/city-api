package users

type User struct {
	Id          string     `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Gender      Gender     `json:"gender"`
	AccessType  AccessType `json:"access_type"`
	Avatar      string     `json:"avatar"`
}

type UserUpdate struct {
	Id          string      `json:"id"`
	FirstName   *string     `json:"first_name"`
	LastName    *string     `json:"last_name"`
	Username    *string     `json:"username"`
	Password    *string     `json:"password"`
	Email       *string     `json:"email"`
	PhoneNumber *string     `json:"phone_number"`
	Gender      *Gender     `json:"gender"`
	AccessType  *AccessType `json:"access_type"`
	Avatar      *string     `json:"avatar"`
}
type Gender string

var (
	Female             Gender = "female"
	Male               Gender = "male"
	genderTypeToString        = map[Gender]string{
		Female: "female",
		Male:   "male",
	}
	stringToGenderType = map[string]Gender{
		"female": Female,
		"male":   Male,
	}
)

func (c Gender) ToString() string {
	return genderTypeToString[c]
}

func ToGenderType(s string) Gender {
	return stringToGenderType[s]
}

func IsGenderExist(s string) bool {
	genderTypes := []string{"male", "female"}
	for _, v := range genderTypes {
		if v == s {
			return true
		}
	}
	return false
}

type AccessType string

var (
	AccessTypeUser     AccessType = "user"
	AccessTypeAdmin    AccessType = "admin"
	accessTypeToString            = map[AccessType]string{
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
