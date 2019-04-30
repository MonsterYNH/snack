package leavemessage

import (
	"gopkg.in/mgo.v2/bson"
)

type CreateLeaveMessageEntry struct {
	CreateUserName string `json:"create_user_name"`
	CreateUserAvatar string `json:"create_user_avatar"`
	Content string `json:"content"`
	CreateUser *bson.ObjectId `json:"create_user"`
}

type CommentLeaveMessageEntry struct {
	MessageId bson.ObjectId `json:"message_id"`
	Content string	`json:"content"`
	CreateUserName string `json:"create_user_name"`
	CreateUserAvatar string `json:"create_user_avatar"`
	CreateUser *bson.ObjectId `json:"create_user"`
}