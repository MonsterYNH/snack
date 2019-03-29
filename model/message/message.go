package model

import (
	"snack/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID         bson.ObjectId  `json:"-" bson:"_id"`
	UserID     bson.ObjectId  `json:"user_id" bson:"user_id"`
	FromUserID bson.ObjectId  `json:"from_user_id" bson:"from_user_id"`
	Type       int            `json:"type" bson:"type"`
	ObjectId   *bson.ObjectId `json:"-" bson:"object_id"`
	Content    string         `json:"content" bson:"content"`
	Status     int            `json:"-" bson:"status"`
	CreateTime int64          `json:"-" bson:"create_time"`
}

func CreateUserMessage(userID, fromUserID bson.ObjectId, content string, messageType int, objectID *bson.ObjectId) (*Message, error) {
	session := db.GetMgoSession()
	defer session.Close()

	message := Message{
		ID:         bson.NewObjectId(),
		UserID:     userID,
		FromUserID: fromUserID,
		Type:       messageType,
		ObjectId:   objectID,
		Content:    content,
		Status:     0,
		CreateTime: time.Now().Unix(),
	}
	err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Insert(message)
	return &message, err
}

func GetUserMessage(userID bson.ObjectId) ([]Message, error) {
	session := db.GetMgoSession()
	defer session.Close()

	messages := make([]Message, 0)
	err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Find(bson.M{"user_id": userID}).All(&messages)
	return messages, err
}

func GetUserMessageCount(userID bson.ObjectId) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Find(bson.M{"user_id": userID}).Count()
}
