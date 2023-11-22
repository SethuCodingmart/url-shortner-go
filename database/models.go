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

type Workspaces struct {
	Id            int        `json:"id"`
	Name          string     `json:"name"`
	Shorthandname string     `json:"shorthandname"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"createdat"`
	UpdatedAt     time.Time  `json:"updatedat"`
	DeletedAT     *time.Time `json:"deletedat"`
}

type Transcation_for string

const (
	MAIL        Transcation_for = "MAIL"
	PHONE       Transcation_for = "PHONE"
	SHORTEN     Transcation_for = "SHORTEN"
	PROFILE_BIO Transcation_for = "PROFILE_BIO"
	ALL         Transcation_for = "ALL"
)

type Transcation_type string

const (
	CREDITED   Transcation_type = "CREDITED"
	DEBITED    Transcation_type = "DEBITED"
	DECLINED   Transcation_type = "DECLINED"
	PENDING    Transcation_type = "PENDING"
	PROCESSING Transcation_type = "PROCESSING"
)

type Transactions struct {
	Id        int              `json:"id"`
	UserID    int              `json:"user_id"`
	Tfor      Transcation_for  `json:"tfor"`
	Type      Transcation_type `json:"type"`
	Value     int              `json:"value"`
	CreatedAt time.Time        `json:"createdat"`
	UpdatedAt *time.Time       `json:"updatedat,omitempty"`
	DeletedAT *time.Time       `json:"deletedat,omitempty"`
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
