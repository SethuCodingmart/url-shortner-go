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

type Register struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type SendSignUpOTP struct {
	Gmail string `json:"gmail"`
}

type ForgetPassword struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type OTPType string

const (
	FORGOT_PASSWORD OTPType = "FORGOT_PASSWORD"
	SIGNUP          OTPType = "SIGNUP"
	LOGIN           OTPType = "LOGIN"
)

type SaveOTP struct {
	Key   string  `json:"key"`
	Value string  `json:"value"`
	Type  OTPType `json:"type"`
}
