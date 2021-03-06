package model

import (
	"errors"
	"fmt"
	"snack/db"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	// 基本信息
	ID            bson.ObjectId `json:"id" bson:"_id"`
	Account       string        `json:"account" bson:"account"`
	Password      string        `json:"-" bson:"password"`
	Name          string        `json:"name" bson:"name"`
	Sex           int           `json:"sex" bson:"sex"`
	Age           int           `json:"age" bson:"age"`
	Avatar        string        `json:"avatar" bson:"avatar"`
	SignTure      string        `json:"signture" bson:"signture"`
	City          string        `json:"city" bson:"city"`
	Country       string        `json:"country" bson:"country"`
	Phone         string        `json:"phone" bson:"phone"`
	Email         string        `json:"email" bson:"email"`
	Status        int           `json:"-" bson:"status"`
	CreateTime    int64         `json:"-" bson:"create_time"`
	LastLoginTime int64         `json:"-" bson:"last_login_time"`

	// 用户关注信息
	UserFollowed    []bson.ObjectId `json:"-" bson:"user_followed"`
	UserFollowedNum int             `json:"user_followed_num" bson:"user_followed_num"`
	FollowedUser    []bson.ObjectId `json:"-" bson:"followed_user"`
	FollowedUserNum int             `json:"followed_user_num" bson:"followed_user_num"`
	FollowStatus    string          `json:"follow_status" bson:"follow_status"`
}

const (
	USER_STATUS_UNFOLLOW    = "unfollow"
	USER_STATUS_FOLLOW      = "follow"
	USER_STATUS_FOLLOW_EACH = "each"
)

func CreateUser(user User) (*User, error) {
	session := db.GetMgoSession()
	defer session.Close()

	user.ID = bson.NewObjectId()
	user.Status = db.STATUS_USER_NORMAL
	user.CreateTime = time.Now().Unix()
	user.LastLoginTime = time.Now().Unix()
	if len(user.Name) == 0 {
		user.Name = strconv.FormatInt(time.Now().Unix(), 10)
	}
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Insert(user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers(query bson.M, start, limit int, args ...string) ([]User, error) {
	session := db.GetMgoSession()
	defer session.Close()

	users := make([]User, 0)
	var err error
	project := make(map[string]int)
	if len(args) > 0 {
		for _, arg := range args {
			project[arg] = 1
		}
		err = session.DB(db.DB_NAME).C(db.COLL_USER).Find(query).Select(project).Skip((start - 1) * limit).Limit(limit).All(&users)
	} else {
		err = session.DB(db.DB_NAME).C(db.COLL_USER).Find(query).Skip((start - 1) * limit).Limit(limit).All(&users)
	}

	return users, err
}

func GetUser(query bson.M, args ...string) (*User, error) {
	session := db.GetMgoSession()
	defer session.Close()

	var userEntry User
	var err error
	project := make(map[string]int)
	if len(args) > 0 {
		for _, arg := range args {
			project[arg] = 1
		}
		err = session.DB(db.DB_NAME).C(db.COLL_USER).Find(query).Select(project).One(&userEntry)
	} else {
		err = session.DB(db.DB_NAME).C(db.COLL_USER).Find(query).One(&userEntry)
	}

	return &userEntry, err
}

func GetUserCount(query bson.M) (int, error) {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_USER).Find(query).Count()
}

func GetUserById(id string, args ...string) (*User, error) {
	session := db.GetMgoSession()
	defer session.Close()

	var userEntry User
	var err error
	project := make(map[string]int)
	if len(args) > 0 {
		for _, arg := range args {
			project[arg] = 1
		}
		err = session.DB(db.DB_NAME).C(db.COLL_USER).FindId(bson.ObjectIdHex(id)).Select(project).One(&userEntry)
	} else {
		err = session.DB(db.DB_NAME).C(db.COLL_USER).FindId(bson.ObjectIdHex(id)).One(&userEntry)
	}

	return &userEntry, err
}

func (user *User) UpdateUser(query, update bson.M) error {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_USER).Update(query, update)
}

func (user *User) UpdateAllUser(query, update bson.M) (*mgo.ChangeInfo, error) {
	session := db.GetMgoSession()
	defer session.Close()

	return session.DB(db.DB_NAME).C(db.COLL_USER).UpdateAll(query, update)
}

func FollowUser(option string, userId, operatorId bson.ObjectId) error {
	session := db.GetMgoSession()
	defer session.Close()

	switch option {
	case "follow":
		if err := session.DB(db.DB_NAME).C(db.COLL_USER).Update(bson.M{"_id": operatorId, "followed_user": bson.M{"$ne": userId}}, bson.M{"$addToSet": bson.M{"followed_user": userId}, "$inc": bson.M{"followed_user_num": 1}}); err != nil {
			return err
		}
		if err := session.DB(db.DB_NAME).C(db.COLL_USER).Update(bson.M{"_id": userId, "user_followed": bson.M{"$ne": operatorId}}, bson.M{"$addToSet": bson.M{"user_followed": operatorId}, "$inc": bson.M{"user_followed_num": 1}}); err != nil {
			return err
		}
	case "unfollow":
		if err := session.DB(db.DB_NAME).C(db.COLL_USER).Update(bson.M{"_id": operatorId, "followed_user": userId}, bson.M{"$pull": bson.M{"followed_user": userId}, "$inc": bson.M{"followed_user_num": -1}}); err != nil {
			return err
		}
		if err := session.DB(db.DB_NAME).C(db.COLL_USER).Update(bson.M{"_id": userId, "user_followed": operatorId}, bson.M{"$pull": bson.M{"user_followed": operatorId}, "$inc": bson.M{"user_followed_num": -1}}); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("not support option %s", option))
	}

	return nil
}

func (user *User) GetUserFollowStatus(operatorId string) error {
	session := db.GetMgoSession()
	defer session.Close()

	userEntry := User{}
	user.FollowStatus = USER_STATUS_UNFOLLOW
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Find(bson.M{"_id": bson.ObjectIdHex(operatorId), "followed_user": user.ID}).One(&userEntry); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		} else {
			return err
		}
	}
	user.FollowStatus = USER_STATUS_FOLLOW
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Find(bson.M{"_id": user.ID, "followed_user": operatorId}).One(&userEntry); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		} else {
			return err
		}
	}
	user.FollowStatus = USER_STATUS_FOLLOW_EACH
	return nil
}

func GetUserFollowed(id bson.ObjectId, operatorId *bson.ObjectId, kind string, start, limit int) (map[string]interface{}, error) {
	session := db.GetMgoSession()
	defer session.Close()

	if kind != "user_followed" && kind != "followed_user" {
		return nil, errors.New("wrong parameter kind")
	}

	aggregate := []bson.M{
		bson.M{"$match": bson.M{"_id": id}},
		bson.M{"$unwind": "$" + kind},
		bson.M{"$skip": (start - 1) * limit},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{"from": "user", "localField": kind, "foreignField": "_id", "as": "user_info"}},
		bson.M{"$unwind": "$user_info"},
		bson.M{"$project": bson.M{
			"_id":               "$user_info._id",
			"name":              "$user_info.name",
			"sex":               "$user_info.sex",
			"avatar":            "$user_info.avatar",
			"signture":          "$user_info.signture",
			"user_followed_num": "$user_info.user_followed_num",
			"followed_user_num": "$user_info.followed_user_num",
		}},
	}

	users := make([]User, 0)
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Pipe(aggregate).All(&users); err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make(map[string]interface{})
	result["data"] = users

	var total map[string]interface{}
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Pipe([]bson.M{
		bson.M{"$match": bson.M{"_id": id}},
		bson.M{"$unwind": "$" + kind},
		bson.M{"$group": bson.M{
			"_id":   "$" + kind,
			"total": bson.M{"$sum": 1},
		}},
	}).One(&total); err != nil {
		if err == mgo.ErrNotFound {
			result["total"] = 0
			return result, nil
		} else {
			return nil, err
		}
	}

	result["total"] = total["total"]
	return result, nil

}
