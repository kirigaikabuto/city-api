package application

type Application struct {
	Id          string      `json:"id"`
	Address     string      `json:"address"`
	AppType     ProblemType `json:"app_type"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Patronymic  string      `json:"patronymic"`
	PhoneNumber string      `json:"phone_number"`
	PhotoUrl    string      `json:"photo_url"`
	VideoUrl    string      `json:"video_url"`
	CreatedDate string      `json:"created_date"`
}

type ProblemType string

var (
	Dump                  ProblemType = "свалка"
	OversizeWaste         ProblemType = "крупногабаритные отходы"
	OverflowingContainers ProblemType = "переполненные контейнеры"
	OverflowingBins       ProblemType = "переполненные урны"
)

var (
	problemTypeToString = map[ProblemType]string{
		Dump:                  "свалка",
		OversizeWaste:         "крупногабаритные отходы",
		OverflowingContainers: "переполненные контейнеры",
		OverflowingBins:       "переполненные урны",
	}
	stringToProblemType = map[string]ProblemType{
		"свалка":                   Dump,
		"крупногабаритные отходы":  OversizeWaste,
		"переполненные контейнеры": OverflowingContainers,
		"переполненные урны":       OverflowingBins,
	}
)

func (c ProblemType) ToString() string {
	return problemTypeToString[c]
}

func ToProblemType(s string) ProblemType {
	return stringToProblemType[s]
}

func IsProblemTypeExist(s string) bool {
	problemTypes := []string{"свалка", "крупногабаритные отходы", "переполненные контейнеры", "переполненные урны"}
	for _, v := range problemTypes {
		if v == s {
			return true
		}
	}
	return false
}
