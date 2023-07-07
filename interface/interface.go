package interfaceGo

type URLParameters struct {
	Id       int    `json:"id"`
	Location string `json:"location"`
	Alias    string `json:"alias"`
}

type URLResultSuccess struct {
	Alias   string `json:"alias"`
	Success bool   `json:"success"`
}

type HealthCheck struct {
	DB  string `json: "db"`
	APP string `json:"app"`
}

type Login struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}
