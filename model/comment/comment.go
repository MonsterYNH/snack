package comment

import (
	"snack/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	ID               bson.ObjectId   `json:"id" bson:"_id"`
	Type             string          `json:"-" bson:"type"`
	ObjectId         bson.ObjectId   `json:"object_id" bson:"object_id"`
	Content          string          `json:"content" bson:"content"`
	Comments         []Reply         `json:"comments" bson:"comments"`
	CommentNum       int             `json:"comment_num" bson:"comment_num"`
	LikeNum          int             `json:"like_num" bson:"like_num"`
	LikeUsers        []bson.ObjectId `json:"-" bson:"like_users"`
	CreateUser       bson.ObjectId   `json:"create_user" bson:"create_user"`
	CreateUserName   string          `json:"create_user_name" bson:"create_user_name"`
	CreateUserAvatar string          `json:"create_user_avatar" bson:"create_user_avatar"`
	CreateTime       int64           `json:"create_time" bson:"create_time"`
	Status           int             `json:"-" bson:"status"`
}

type Reply struct {
	ID           bson.ObjectId   `json:"id" bson:"_id"`
	ToUser       bson.ObjectId   `json:"to_user" bson:"to_user"`
	ToUserAvatar string          `json:"to_user_avatar" bson:"to_user_avatar"`
	ToUserName   string          `json:"to_user_name" bson:"to_user_name"`
	User         bson.ObjectId   `json:"user" bson:"user"`
	UserAvatar   string          `json:"user_avatar" bson:"user_avatar"`
	UserName     string          `json:"user_name" bson:"user_name"`
	Content      string          `json:"content" bson:"content"`
	Comments     []Reply         `json:"comments" bson:"comments"`
	CommentNum   int             `json:"comment_num" bson:"comment_num"`
	LikeNum      int             `json:"like_num" bson:"like_num"`
	LikeUsers    []bson.ObjectId `json:"-" bson:"like_users"`
	CreateTime   int64           `json:"create_time" bson:"create_time"`
	Status       int             `json:"-" bson:"status"`
}

type ReplyEntry struct {
	Replies []Reply `json:"comments" bson:"comments"`
}

func GetComments(query bson.M, start, limit int) ([]Comment, error) {
	session := db.GetMgoSession()
	defer session.Close()

	comments := make([]Comment, 0)
	err := session.DB(db.DB_NAME).C(db.COLL_COMMENT).Find(query).Skip((start - 1) * limit).Limit(start * limit).All(&comments)
	return comments, err
}

func GetCommentsByPage(kind string, objectId bson.ObjectId, userId *bson.ObjectId, start, limit, replyLimit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	comments := make([]Comment, 0)
	result := make(map[string]interface{})
	match := bson.M{"status": 0, "type": kind, "object_id": objectId}
	if err := session.DB(db.DB_NAME).C(db.COLL_COMMENT).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{
			"from":         "user",
			"localField":   "create_user",
			"foreignField": "_id",
			"as":           "owner",
		}},
		bson.M{"$unwind": "$owner"},
		bson.M{"$project": bson.M{
			"object_id":          1,
			"content":            1,
			"create_user":        1,
			"comments":           1,
			"comment_num":        1,
			"like_num":           1,
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userId}},
			"create_user_name":   "$owner.name",
			"create_user_avatar": "$owner.avatar",
			"create_time":        1,
		}},
	}).All(&comments); err != nil {
		return result, err
	}
	total, err := session.DB(db.DB_NAME).C(db.COLL_COMMENT).Find(match).Count()
	if err != nil {
		return result, err
	}
	result["data"] = comments
	result["total"] = total
	return result, err
}

func GetCommentNum(kind string, objectId bson.ObjectId) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	match := bson.M{"status": 0, "type": kind, "object_id": objectId}
	return session.DB(db.DB_NAME).C(db.COLL_COMMENT).Find(match).Count()
}

func GetCommentRepliesByPage(commentId bson.ObjectId, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	replies := make([]Reply, 0)
	result := make(map[string]interface{})
	match := bson.M{"_id": commentId, "status": 0}
	err := session.DB(db.DB_NAME).C(db.COLL_COMMENT).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$unwind": "$comments"},
		bson.M{"$sort": bson.M{"comments.create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
		bson.M{"$project": bson.M{"_id": 0, "comments": 1}},
		bson.M{"$unwind": "$comments"},
		bson.M{"$lookup": bson.M{
			"from":         "user",
			"localField":   "comments.to_user",
			"foreignField": "_id",
			"as":           "create_user",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "user",
			"localField":   "comments.user",
			"foreignField": "_id",
			"as":           "comment_user",
		}},
		bson.M{"$unwind": "$comment_user"},
		bson.M{"$unwind": "$create_user"},
		bson.M{"$project": bson.M{
			"_id":            "$comments._id",
			"to_user":        "$comments.to_user",
			"to_user_avatar": "$comment_user.avatar",
			"to_user_name":   "$comment_user.name",
			"user":           "$comments.user",
			"user_avatar":    "$create_user.avatar",
			"user_name":      "$create_user.name",
			"content":        "$comments.content",
			"like_num":       "$comments.like_num",
			"create_time":    "$comments.create_time",
		}},
	}).All(&replies)
	comment := Comment{}
	if err := session.DB(db.DB_NAME).C(db.COLL_COMMENT).Find(match).One(&comment); err != nil {
		return result, err
	}
	total := comment.CommentNum
	result["data"] = replies
	result["total"] = total
	return result, err
}

func CreateComment(kind, content string, objectId, createUser bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	comment := Comment{
		ID:         bson.NewObjectId(),
		Type:       kind,
		ObjectId:   objectId,
		CreateUser: createUser,
		Content:    content,
		CreateTime: time.Now().Unix(),
	}

	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).UpdateId(objectId, bson.M{"$inc": bson.M{"comment_num": 1}}); err != nil {
		return err
	}

	return session.DB(db.DB_NAME).C(db.COLL_COMMENT).Insert(comment)
}

func CreateCommentReply(commentId, toUser, user bson.ObjectId, content string) error {
	session := db.GetMgoSession()
	defer session.Close()

	reply := Reply{
		ID:         bson.NewObjectId(),
		ToUser:     toUser,
		User:       user,
		Content:    content,
		CreateTime: time.Now().Unix(),
	}

	return session.DB(db.DB_NAME).C(db.COLL_COMMENT).UpdateId(commentId, bson.M{"$addToSet": bson.M{"comments": reply}, "$inc": bson.M{"comment_num": 1}})
}

func LikeComment(commentId, userID bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_COMMENT).UpdateId(commentId, bson.M{"$inc": bson.M{"like_num": 1}, "$addToSet": bson.M{"like_users": userID}})
}

func LikeCommentReply(commentId, replyId, userId bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_COMMENT).Update(bson.M{"_id": commentId, "comments._id": replyId}, bson.M{"$inc": bson.M{"comments.$.like_num": 1}, "$addToSet": bson.M{"like_users": userId}})
}
