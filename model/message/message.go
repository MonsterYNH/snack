package model

import (
	"errors"
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
	CreateTime     int64         `json:"time" bson:"time"`
}

const (
	// 系统消息
	SYSTEM_MESSAGE = 0
	// 1XXX 系统信息保留字段

	// 2XXX 文章相关
	USER_ARTICLE_LIKE          = 2001
	USER_ARTICLE_COLLECT       = 2002
	USER_ARTICLE_COMMENT       = 2003
	USER_ARTICLE_COMMENT_LIKE  = 2004
	USER_ARTICLE_COMMENT_REPLY = 2005
	// 3XXX 关注相关
	USER_FOLLOW       = 3001
	USER_FOLLOW_TOPIC = 3002

	// 4XXX 随笔相关
)

func SaveMessage(userID, fromUserID bson.ObjectId, content string, messageType int, objectID *bson.ObjectId) (*Message, error) {
	switch messageType {
	case SYSTEM_MESSAGE:
		return CreateUserMessage(userID, fromUserID, "系统消息", content, messageType, objectID)
	case USER_ARTICLE_LIKE:
		return CreateUserMessage(userID, fromUserID, "点赞了你的文章", content, messageType, objectID)
	case USER_ARTICLE_COLLECT:
		return CreateUserMessage(userID, fromUserID, "收藏了你的文章", content, messageType, objectID)
	case USER_ARTICLE_COMMENT:
		return CreateUserMessage(userID, fromUserID, "评论了你的文章", content, messageType, objectID)
	case USER_ARTICLE_COMMENT_LIKE:
		return CreateUserMessage(userID, fromUserID, "在文章评论中点赞了你", content, messageType, objectID)
	case USER_ARTICLE_COMMENT_REPLY:
		return CreateUserMessage(userID, fromUserID, "在文章评论中回复了你", content, messageType, objectID)
	case USER_FOLLOW:
		return CreateUserMessage(userID, fromUserID, "关注了你", content, messageType, objectID)
	default:
		return nil, errors.New("not match any message type")
	}
}

func CreateUserMessage(userID, fromUserID bson.ObjectId, title, content string, messageType int, objectID *bson.ObjectId) (*Message, error) {
	session := db.GetMgoSession()
	defer session.Close()

	message := Message{
		ID:         bson.NewObjectId(),
		UserID:     userID,
		FromUserID: fromUserID,
		Type:       messageType,
		ObjectId:   objectID,
		Content:    content,
		Title:      title,
		Status:     0,
		CreateTime: time.Now().Unix(),
	}
	err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Insert(message)
	return &message, err
}

func GetUserMessage(userID string, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	messageInfos := make([]MessageInfo, 0)
	result := make(map[string]interface{})
	match := bson.M{"user_id": bson.ObjectIdHex(userID), "status": 0}
	aggregate := []bson.M{
		bson.M{"$match": match},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "from_user_id", "foreignField": "_id", "as": "user_info"}},
		bson.M{"$unwind": "$user_info"},
		bson.M{"$project": bson.M{"from_user_id": "$from_user_id", "from_user_avatar": "$user_info.avatar", "object_id": "$user_info.object_id", "title": "$title", "content": "$content", "type": "$type", "id": "$_id", "_id": 0, "time": "$create_time"}},
	}
	total, err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Find(match).Count()
	if err != nil {
		return result, err
	}
	err = session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Pipe(aggregate).All(&messageInfos)
	result["total"] = total
	result["data"] = messageInfos
	return result, err
}

func GetUserMessageCount(userID string) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Find(bson.M{"user_id": bson.ObjectIdHex(userID), "status": 0}).Count()
}

func SetUserMessageRead(messageID bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_MESSAGE).UpdateId(messageID, bson.M{"$set": bson.M{"status": 1}})
}

func SetUserMessageReadAll(userId bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	_, err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).UpdateAll(bson.M{"user_id": userId}, bson.M{"$set": bson.M{"status": 1}})
	return err
}
