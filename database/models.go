package database

import (
	"time"
	interfaceGo "urlShortner/interface"
)

type Users struct {
	Id        int        `json:"id"`
	Gmail     string     `json:"gmail"`
	Password  string     `json:"password"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	DeletedAt *time.Time `json:"deletedat"`
	CreatedAt time.Time  `json:"createdat"`
	UpdatedAt time.Time  `json:"updatedat"`
}

type OTP struct {
	Id        int                 `json:"id"`
	Key       string              `json:"key"`
	Value     string              `json:"value"`
	Type      interfaceGo.OTPType `json:"type"`
	CreatedAt time.Time           `json:"createdat"`
	UpdatedAt time.Time           `json:"updatedat"`
	DeletedAT *time.Time          `json:"deletedat"`
}

type Urls struct {
	Id        int        `json:"id"`
	Location  string     `json:"location"`
	Alias     string     `json:"alias"`
	Expiresat *time.Time `json:"expiresat"`
	IsExpired bool       `json:"isexpired"`
	UserId    int        `json:"userid"`
	Createdat time.Time  `json:"createdat"`
}
