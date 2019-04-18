package article

import "gopkg.in/mgo.v2/bson"

type ArticleCreate struct {
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	Images     []string      `json:"images"`
	CreateUser bson.ObjectId `json:"create_user"`
}

type ArticleListEntry struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}
