package model

import (
	"snack/db"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
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
}

func CreateUser(user *User) (*User, error) {
	session := db.GetMgoSession()
	defer session.Close()

	user.ID = bson.NewObjectId()
	user.Status = db.STATUS_USER_NORMAL
	user.CreateTime = time.Now().Unix()
	user.LastLoginTime = time.Now().Unix()
	if err := session.DB(db.DB_NAME).C(db.COLL_USER).Insert(user); err != nil {
		return nil, err
	}
	return user, nil
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
