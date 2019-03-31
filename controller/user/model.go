package user

import (
	model "snack/model/user"
)

type UserLoginEntry struct {
	Type     string `json:"type" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password" binding:"required"`
}

type UserRegistEntry struct {
	// model.User
	// Agreement bool   `json:"agreement"`
	Type     string `json:"type" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
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
