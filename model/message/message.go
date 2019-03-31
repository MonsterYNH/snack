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
	Title      string         `json:"title" bson:"title"`
	Status     int            `json:"-" bson:"status"`
	CreateTime int64          `json:"-" bson:"create_time"`
}

type MessageInfo struct {
	ID             bson.ObjectId `json:"id" bson:"id"`
	FromUserID     bson.ObjectId `json:"from_user_id" bson:"from_user_id"`
	FromUserAvatar string        `json:"from_user_avatar" bson:"from_user_avatar"`
	ObjectID       bson.ObjectId `json:"object_id" bson:"object_id"`
	Content        string        `json:"content" bson:"content"`
	Title          string        `json:"title" bson:"title"`
	Type           int           `json:"type" bson:"type"`
}

func CreateUserMessage(userID, fromUserID bson.ObjectId, tiitle, content string, messageType int, objectID *bson.ObjectId) (*Message, error) {
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

func GetUserMessage(userID bson.ObjectId) ([]MessageInfo, error) {
	session := db.GetMgoSession()
	defer session.Close()

	messageInfos := make([]MessageInfo, 0)
	aggregate := []bson.M{
		bson.M{"$match": bson.M{"user_id": userID, "status": 0}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "from_user_id", "foreignField": "_id", "as": "user_info"}},
		bson.M{"$unwind": "$user_info"},
		bson.M{"$project": bson.M{"from_user_id": "$from_user_id", "from_user_avatar": "$user_info.avatar", "object_id": "$user_info.object_id", "title": "$title", "content": "$content", "type": "$type", "id": "$_id", "_id": 0}},
	}
	err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Pipe(aggregate).All(&messageInfos)
	return messageInfos, err
}

func GetUserMessageCount(userID bson.ObjectId) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Find(bson.M{"user_id": userID}).Count()
}

func SetUserMessageRead(messageID bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_MESSAGE).UpdateId(messageID, bson.M{"$set": bson.M{"status": 1}})
}
