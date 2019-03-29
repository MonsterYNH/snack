package user

import (
	model "snack/model/user"
)

type UserLoginEntry struct {
	Type     string `json:"type"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserRegistEntry struct {
	model.User
	Agreement bool   `json:"agreement"`
	Confirm   string `json:"confirm"`
}

type UserInfoEntry struct {
	model.User
	MessageCount int `json:"message_count"`
}

type UserPageEntry struct {
	Start string `json:"start"`
	Limit string `json:"limit"`
}

type UserPage struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}
