package model

import (
	"time"
)

type Account struct {
	AccountID int64     `json:"account_id"`
	Email     string    `json:"user_email"`
	Fullname  string    `json:"user_fullname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Accounts []*Account

func NewAccount(email, fullname string) *Account {
	return &Account{
		Email:    email,
		Fullname: fullname,
	}
}
