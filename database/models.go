package database

import (
	"time"
)

type Urls struct {
	Id        int       `json:"id"`
	Location  string    `json:"location"`
	Alias     string    `json:"alias"`
	Expiresat time.Time `json:"expiresat"`
	IsExpired bool      `json:"isexpired"`
	UserId    int       `json:"userid"`
	Createdat time.Time `json:"createdat"`
}

type Users struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Gmail     string    `json:"gmail"`
	Password  string    `json:"password"`
	Authkey   string    `json:"authkey"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
