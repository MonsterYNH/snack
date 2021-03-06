package article

import (
	"errors"
	"snack/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	ID              bson.ObjectId   `json:"id" bson:"_id"`
	Tags            []string        `json:"tags" bson:"tags"`
	Category        string          `json:"category" bson:"category"`
	Images          []string        `json:"images" bson:"images"`
	Title           string          `json:"title" bson:"title"`
	Description     string          `json:"description" bson:"description"`
	Content         string          `json:"content" bson:"content"`
	CreateUser      bson.ObjectId   `json:"create_user" bson:"create_user"`
	LikeUsers       []bson.ObjectId `json:"-" bson:"like_users"`
	LikeUsersNum    int             `json:"like_users_num" bson:"like_users_num"`
	CollectUsers    []bson.ObjectId `json:"-" bson:"collect_users"`
	CollectUsersNum int             `json:"collect_users_num" bson:"collect_users_num"`
	CommentNum      int             `json:"comment_num" bson:"comment_num"`
	ReadNum         int             `json:"read_num" bson:"read_num"`
	CreateTime      int64           `json:"create_time" bson:"create_time"`
	UpdateTime      int64           `json:"update_time" bson:"update_time"`
	Weight          int             `json:"weight" bson:"weight"`
	Status          int             `json:"-" bson:"status"`
}

type ArticleEntry struct {
	ID               bson.ObjectId   `json:"id" bson:"_id"`
	Tags             []string        `json:"tags" bson:"tags"`
	Images           []string        `json:"images" bson:"images"`
	Title            string          `json:"title" bson:"title"`
	Content          string          `json:"content" bson:"content"`
	Description      string          `json:"description" bson:"description"`
	CreateUser       bson.ObjectId   `json:"create_user" bson:"create_user"`
	LikeUsers        []bson.ObjectId `json:"-" bson:"like_users"`
	LikeUsersNum     int             `json:"like_users_num" bson:"like_users_num"`
	CollectUsers     []bson.ObjectId `json:"-" bson:"collect_users"`
	CollectUsersNum  int             `json:"collect_users_num" bson:"collect_users_num"`
	CommentNum       int             `json:"comment_num" bson:"comment_num"`
	CreateTime       int64           `json:"create_time" bson:"create_time"`
	UpdateTime       int64           `json:"update_time" bson:"update_time"`
	ReadNum          int             `json:"read_num" bson:"read_num"`
	Status           int             `json:"-" bson:"status"`
	LikeStatus       int             `json:"like_status" bson:"like_status"`
	CollectStatus    int             `json:"collect_status" bson:"collect_status"`
	CreateUserName   string          `json:"create_user_name" bson:"create_user_name"`
	CreateUserAvatar string          `json:"create_user_avatar" bson:"create_user_avatar"`
}

type TagEntry struct {
	Tag   string `json:"tag" bson:"tag"`
	Total int    `json:"total" bson:"total"`
}

type CategoryEntry struct {
	Category string `json:"category" bson:"category"`
	Total    int    `json:"total" bson:"total"`
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

func GetArticleListByTag(tags []string, userID *bson.ObjectId, start, limit, sort int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	result := make(map[string]interface{})
	articles := make([]ArticleEntry, 0)
	match := bson.M{}
	if len(tags) == 0 {
		match = bson.M{"status": 0}
	} else {
		match = bson.M{"status": 0, "tags": bson.M{"$in": tags}}
	}
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userID}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userID}},
		}},
		bson.M{"$sort": bson.M{"create_time": sort}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles)
	if err != nil {
		return nil, err
	}
	if total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count(); err != nil {
		return nil, err
	} else {
		result["data"] = articles
		result["total"] = total
	}
	return result, err
}

func GetArticleByUserId(userId bson.ObjectId, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]ArticleEntry, 0)
	result := make(map[string]interface{})
	match := bson.M{"create_user": userId, "status": 0}
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userId}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userId}},
		}},
		bson.M{"$sort": bson.M{"create_time": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles); err != nil {
		return nil, err
	}
	if total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count(); err != nil {
		result["total"] = total
		result["data"] = articles
		return nil, err
	} else {
		result["total"] = total
		result["data"] = articles
		return result, err
	}
}

func GetHotArticle(userID *bson.ObjectId, start, limit, sort int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	result := make(map[string]interface{})
	articles := make([]ArticleEntry, 0)
	match := bson.M{"status": 0}
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userID}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userID}},
		}},
		bson.M{"$sort": bson.M{"read_num": sort}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles)
	if err != nil {
		return nil, err
	}
	total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count()
	if err != nil {
		return nil, err
	}
	result["data"] = articles
	result["total"] = total
	return result, err
}

func GetNewArticle(userID *bson.ObjectId, start, limit, sort int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	result := make(map[string]interface{})
	articles := make([]ArticleEntry, 0)
	match := bson.M{"status": 0}
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userID}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userID}},
		}},
		bson.M{"$sort": bson.M{"create_time": sort}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles)
	if err != nil {
		return nil, err
	}
	total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count()
	if err != nil {
		return nil, err
	}
	result["data"] = articles
	result["total"] = total
	return result, err
}

func GetRecommendArticle(userID *bson.ObjectId, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	result := make(map[string]interface{})
	articles := make([]ArticleEntry, 0)
	match := bson.M{"status": 0}
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe([]bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userID}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userID}},
		}},
		bson.M{"$sort": bson.M{"weight": -1}},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": start * limit},
	}).All(&articles)
	if err != nil {
		return nil, err
	}
	total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count()
	if err != nil {
		return nil, err
	}
	result["data"] = articles
	result["total"] = total
	return result, err
}

func GetSameArticle(id bson.ObjectId, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	articles := make([]ArticleEntry, 0)
	article := ArticleEntry{}
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(bson.M{"_id": id, "status": 0}).One(&article); err != nil {
		return nil, err
	}
	match := bson.M{"status": 0, "tags": bson.M{"$in": article.Tags}}
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Skip((start - 1) * limit).Limit(start * limit).All(&articles); err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	total, err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count()
	if err != nil {
		return nil, err
	}
	result["data"] = articles
	result["total"] = total
	return result, nil
}

func GetArticleCountByTag(tag string) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	match := bson.M{"status": 0, "tags": bson.M{"$in": []string{tag}}}
	return session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Find(match).Count()
}

func GetArticleById(id bson.ObjectId, userID *bson.ObjectId) (*ArticleEntry, error) {
	session := db.GetMgoSession()
	defer session.Close()

	article := ArticleEntry{}
	match := bson.M{"_id": id, "status": 0}
	aggregate := []bson.M{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "create_user", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"tags":               "$tags",
			"images":             "$images",
			"title":              "$title",
			"content":            "$content",
			"description":        "$description",
			"like_users_num":     "$like_users_num",
			"collect_users_num":  "$collect_users_num",
			"create_time":        1,
			"update_time":        1,
			"create_user":        "$create_user",
			"create_user_name":   "$user.name",
			"create_user_avatar": "$user.avatar",
			"comment_num":        "$comment_num",
			"read_num":           "$read_num",
			"like_status":        bson.M{"$indexOfArray": []interface{}{"$like_users", userID}},
			"collect_status":     bson.M{"$indexOfArray": []interface{}{"$collect_users", userID}},
		}},
	}
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe(aggregate).One(&article); err != nil {
		return nil, err
	}
	err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Update(match, bson.M{"$inc": bson.M{"read_num": 1}})
	return &article, err
}

func LikeArticle(id, userID bson.ObjectId, option int) (int, int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	var (
		err          error
		articleEntry ArticleEntry
		num          int
		status       int
	)

	if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).FindId(id).One(&articleEntry); err != nil {
		return 0, -1, err
	}

	switch option {
	case 1:
		// 检查是否点赞
		for _, user := range articleEntry.LikeUsers {
			if user.Hex() == userID.Hex() {
				return 0, -1, errors.New("has liked")
			}
		}
		if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).UpdateId(id, bson.M{"$inc": bson.M{"like_users_num": 1}, "$addToSet": bson.M{"like_users": userID}}); err != nil {
			return 0, -1, err
		}
		status = 1
	case -1:
		// 检查是否
		exist := false
		for _, user := range articleEntry.LikeUsers {
			if user.Hex() == userID.Hex() {
				exist = true
				break
			}
		}
		if !exist {
			return 0, -1, errors.New("please collect")
		}
		if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).UpdateId(id, bson.M{"$inc": bson.M{"like_users_num": -1}, "$pull": bson.M{"like_users": userID}}); err != nil {
			return 0, -1, err
		}
		status = -1
	default:
		return 0, -1, errors.New("option wrong")
	}
	if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).FindId(id).Select(bson.M{"like_users_num": 1}).One(&articleEntry); err != nil {
		return 0, -1, err
	}
	num = articleEntry.LikeUsersNum
	return num, status, err
}

func CollectArticle(id, userID bson.ObjectId, option int) (int, int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	var (
		err          error
		articleEntry ArticleEntry
		num          int
		status       int
	)
	if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).FindId(id).One(&articleEntry); err != nil {
		return 0, -1, err
	}

	switch option {
	case 1:
		// 检查是否收藏
		for _, user := range articleEntry.CollectUsers {
			if user.Hex() == userID.Hex() {
				return 0, -1, errors.New("has collected")
			}
		}
		if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).UpdateId(id, bson.M{"$inc": bson.M{"collect_users_num": 1}, "$addToSet": bson.M{"collect_users": userID}}); err != nil {
			return 0, -1, err
		}
		status = 1
	case -1:
		// 检查是否收藏
		exist := false
		for _, user := range articleEntry.CollectUsers {
			if user.Hex() == userID.Hex() {
				exist = true
				break
			}
		}
		if !exist {
			return 0, -1, errors.New("please collect")
		}
		if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).UpdateId(id, bson.M{"$inc": bson.M{"collect_users_num": -1}, "$pull": bson.M{"collect_users": userID}}); err != nil {
			return 0, -1, err
		}
		status = -1
	default:
		return 0, -1, errors.New("option wrong")
	}
	if err = session.DB(db.DB_NAME).C(db.COLL_ARTICLE).FindId(id).Select(bson.M{"collect_users_num": 1}).One(&articleEntry); err != nil {
		return 0, -1, err
	}
	num = articleEntry.CollectUsersNum
	return num, status, err
}

func GetArticleTags() ([]TagEntry, error) {
	session := db.GetMgoSession()
	defer session.Close()

	aggregate := []bson.M{
		bson.M{"$unwind": "$tags"},
		bson.M{"$group": bson.M{
			"_id":   "$tags",
			"total": bson.M{"$sum": 1},
		}},
		bson.M{"$project": bson.M{
			"tag":   "$_id",
			"total": "$total",
		}},
	}
	tags := make([]TagEntry, 0)
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe(aggregate).All(&tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func GetArticleCategories() ([]CategoryEntry, error) {
	session := db.GetMgoSession()
	defer session.Close()

	aggregate := []bson.M{
		bson.M{"$group": bson.M{
			"_id":   "$category",
			"total": bson.M{"$sum": 1},
		}},
		bson.M{"$project": bson.M{
			"category": "$_id",
			"total":    "$total",
		}},
	}
	categories := make([]CategoryEntry, 0)
	if err := session.DB(db.DB_NAME).C(db.COLL_ARTICLE).Pipe(aggregate).All(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}
