package leavemessage

import (
	"snack/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	LEAVE_MESSAGE_ONLINE  = 0
	LEAVE_MESSAGE_OUTLINE = 1
)

type LeaveMessage struct {
	ID               bson.ObjectId       `json:"id" bson:"_id"`
	Content          string              `json:"content" bson:"content"`
	CreateUser       *bson.ObjectId      `json:"create_user" bson:"create_user"`
	CreateUserName   string              `json:"create_user_name" bson:"create_user_name"`
	CreateUserAvatar string              `json:"create_user_avatar" bson:"create_user_avatar"`
	CreateTime       int64               `json:"create_time" bson:"create_time"`
	Comments         []LeaveMessageEntry `json:"comments" bson:"comments"`
	Status           int                 `json:"status" bson:"status"`
}

type LeaveMessageEntry struct {
	ID               bson.ObjectId       `json:"id" bson:"_id"`
	Content          string              `json:"content" bson:"content"`
	CreateUser       *bson.ObjectId      `json:"create_user" bson:"create_user"`
	CreateUserName   string              `json:"create_user_name" bson:"create_user_name"`
	CreateUserAvatar string              `json:"create_user_avatar" bson:"create_user_avatar"`
	ReplyUserName    string              `json:"reply_user_name" bson:"reply_user_name"`
	ReplyUserAvatar  string              `json:"reply_user_avatar" bson:"reply_user_avatar"`
	CreateTime       int64               `json:"create_time" bson:"create_time"`
	Comments         []LeaveMessageEntry `json:"comments" bson:"comments"`
	Status           int                 `json:"status" bson:"status"`
}

func CreateLeaveMessage(content, createUserName, createUserAvatar string, createUser *bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	leaveMessage := LeaveMessage{
		ID:               bson.NewObjectId(),
		Content:          content,
		CreateTime:       time.Now().Unix(),
		CreateUser:       createUser,
		CreateUserName:   createUserName,
		CreateUserAvatar: createUserAvatar,
		Comments:         nil,
		Status:           LEAVE_MESSAGE_ONLINE,
	}

	return session.DB(db.DB_NAME).C(db.COLL_LEAVE_MESSAGE).Insert(leaveMessage)
}

func CommentLeaveMessage(messageId bson.ObjectId, content, createUserName, createUserAvatar string, createUser *bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	leaveMessage := LeaveMessage{
		ID:               bson.NewObjectId(),
		Content:          content,
		CreateTime:       time.Now().Unix(),
		CreateUser:       createUser,
		CreateUserName:   createUserName,
		CreateUserAvatar: createUserAvatar,
		Comments:         nil,
		Status:           LEAVE_MESSAGE_ONLINE,
	}

	return session.DB(db.DB_NAME).C(db.COLL_LEAVE_MESSAGE).UpdateId(messageId, bson.M{"$addToSet": bson.M{"comments": leaveMessage}})
}

func GetLeaveMessage(start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	leaveMessages := make([]LeaveMessage, 0)
	match := bson.M{"status": LEAVE_MESSAGE_ONLINE}
	aggregate := []bson.M{
		bson.M{"$match": match},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": limit},
	}
	result := make(map[string]interface{})
	if err := session.DB(db.DB_NAME).C(db.COLL_LEAVE_MESSAGE).Pipe(aggregate).All(&leaveMessages); err != nil {
		return result, err
	}
	if num, err := session.DB(db.DB_NAME).C(db.COLL_LEAVE_MESSAGE).Find(match).Count(); err != nil {
		return result, err
	} else {
		result["total"] = num
		result["data"] = leaveMessages
	}
	return result, nil
}
