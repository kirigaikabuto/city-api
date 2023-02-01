package applications

type Application struct {
	Id          string      `json:"id"`
	AppType     ProblemType `json:"app_type"`
	AppStatus   Status      `json:"app_status"`
	Message     string      `json:"message"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Patronymic  string      `json:"patronymic"`
	PhoneNumber string      `json:"phone_number"`
	PhotoUrl    string      `json:"photo_url"`
	VideoUrl    string      `json:"video_url"`
	Address     string      `json:"address"`
	Longitude   float64     `json:"longitude"`
	Latitude    float64     `json:"latitude"`
	UserId      string      `json:"user_id"`
	CreatedDate string      `json:"created_date"`
	Files       []string    `json:"files"`
}

type ApplicationUpdate struct {
	Id          string       `json:"id"`
	PhotoUrl    *string      `json:"photo_url"`
	VideoUrl    *string      `json:"video_url"`
	AppStatus   *Status      `json:"app_status"`
	AppType     *ProblemType `json:"app_type"`
	Message     *string      `json:"message"`
	FirstName   *string      `json:"first_name"`
	LastName    *string      `json:"last_name"`
	Patronymic  *string      `json:"patronymic"`
	PhoneNumber *string      `json:"phone_number"`
	Address     *string      `json:"address"`
	Longitude   *float64     `json:"longitude"`
	Latitude    *float64     `json:"latitude"`
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
		"свалка": Dump,
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

type Status string

var (
	StatusWait        Status = "ожидание"
	StatusOnCheck     Status = "проверка"
	StatusRealization Status = "реализация"
	StatusDone        Status = "выполнен"
)

var (
	statusToString = map[Status]string{
		StatusWait:        "ожидание",
		StatusOnCheck:     "проверка",
		StatusRealization: "реализация",
		StatusDone:        "выполнен",
	}
	stringToStatus = map[string]Status{
		"ожидание":   StatusWait,
		"проверка":   StatusOnCheck,
		"реализация": StatusRealization,
		"выполнен":   StatusDone,
	}
)

func (c Status) ToString() string {
	return statusToString[c]
}

func ToStatus(s string) Status {
	return stringToStatus[s]
}

func IsStatusExist(s string) bool {
	status := []string{"ожидание", "проверка", "реализация", "выполнен"}
	for _, v := range status {
		if v == s {
			return true
		}
	}
	return false
}
