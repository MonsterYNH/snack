package community

import "gopkg.in/mgo.v2/bson"

type CommunityOption struct {
	Kind     string        `json:"kind"`   // 文章  动态 ...
	Type     string        `json:"type"`   // 点赞 收藏
	Option   int           `json:"option"` // 1 点赞 0 取消
	ObjectId bson.ObjectId `json:"object_id"`
}

type CommunityCreate struct {
	Kind     string        `json:"kind"`
	Content  string        `json:"content"`
	ObjectId bson.ObjectId `json:"object_id"`

	CommentUser bson.ObjectId `json:"comment_user"`
}
