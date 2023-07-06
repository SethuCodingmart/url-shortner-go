package interfaceGo

type URLParameters struct {
	Location string `json:"location"`
	Alias    string `json:"alias"`
}

type URLResultSuccess struct {
	Alias   string `json:"alias"`
	Success bool   `json:"success"`
}
