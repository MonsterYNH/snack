package article

import (
	"encoding/json"
	"fmt"
	"snack/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	ID              bson.ObjectId   `json:"id" bson:"_id"`
	Tags            []string        `json:"tags" bson:"tags"`
	Images          []string        `json:"images" bson:"images"`
	Title           string          `json:"title" bson:"title"`
	Content         string          `json:"content" bson:"content"`
	CreateUser      bson.ObjectId   `json:"create_user" bson:"create_user"`
	LikeUsers       []bson.ObjectId `json:"-" bson:"like_users"`
	LikeUsersNum    int             `json:"like_users_num" bson:"like_users_num"`
	DisLikeUsers    []bson.ObjectId `json:"-" bson:"dislike_users"`
	DisLikeUsersNum int             `json:"dislike_users_num" bson:"dislike_users_num"`
	CreateTime      int64           `json:"create_time" bson:"create_time"`
	UpdateTime      int64           `json:"update_time" bson:"update_time"`
	Status          int             `json:"-" bson:"status"`
}

type ArticleEntry struct {
	ID               bson.ObjectId   `json:"id" bson:"_id"`
	Tags             []string        `json:"tags" bson:"tags"`
	Images           []string        `json:"images" bson:"images"`
	Title            string          `json:"title" bson:"title"`
	Content          string          `json:"content" bson:"content"`
	CreateUser       bson.ObjectId   `json:"create_user" bson:"create_user"`
	LikeUsers        []bson.ObjectId `json:"-" bson:"like_users"`
	LikeUsersNum     int             `json:"like_users_num" bson:"like_users_num"`
	DisLikeUsers     []bson.ObjectId `json:"-" bson:"dislike_users"`
	DisLikeUsersNum  int             `json:"dislike_users_num" bson:"dislike_users_num"`
	CreateTime       int64           `json:"create_time" bson:"create_time"`
	UpdateTime       int64           `json:"update_time" bson:"update_time"`
	Status           int             `json:"-" bson:"status"`
	CreateUserName   string          `json:"create_user_name" bson:"create_user_name"`
	CreateUserAvatar string          `json:"create_user_avatar" bson:"create_user_avatar"`
}

func CreateArticle(article Article) error {
	session := db.GetMgoSession()
	defer session.Close()

	now := time.Now().Unix()
	article.ID = bson.NewObjectId()
	article.CreateTime = now
	article.UpdateTime = now
	article.Status = 0

	return session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Insert(article)
}

func GetArticleListByTag(tag string, start, limit int) ([]ArticleEntry, error) {
	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]ArticleEntry, 0)
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": bson.M{"status": 0, "tags": bson.M{"$in": []string{tag}}}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"content":            "$content",
			"like_users_num":     "$like_users_num",
			"dislike_users_num":  "$dislike_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.create_user_name",
			"create_user_avatar": "$user.create_user_avatar",
		}},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles)
	bytes, _ := json.Marshal([]bson.M{
		bson.M{"$match": bson.M{"status": 0, "tags": bson.M{"$in": []string{tag}}}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"content":            "$content",
			"like_users_num":     "$like_users_num",
			"dislike_users_num":  "$dislike_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.create_user_name",
			"create_user_avatar": "$user.create_user_avatar",
		}},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	})
	fmt.Println(string(bytes))
	return articles, err
}

func GetArticleById(id bson.ObjectId) ([]ArticleEntry, error) {
	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]ArticleEntry, 0)
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": bson.M{"status": 1, "_id": id}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$project": bson.M{
			"_id":                0,
			"id":                 "$_id",
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"content":            "$content",
			"like_users_num":     "$like_users_num",
			"dislike_users_num":  "dislike_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$create_user_name",
			"create_user_avatar": "$create_user_avatar",
		}},
	}).All(&articles)
	return articles, err
}
