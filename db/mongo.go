package db

import (
	"dappstatus/global"
	"fmt"
	"reflect"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MongoSession *mgo.Session

const (
	DB_NAME = "blog"

	COLL_USER    = "user"
	COLL_MESSAGE = "message"
	COLL_BANNER  = "banner"
	COLL_ARTICLE = "article"
)

const (
	// 消息状态
	STATUS_MESSAGE_READ   = 1
	STATUS_MESSAGE_UNREAD = 0
	// 用户状态
	STATUS_USER_NORMAL    = 0
	STATUS_USER_ABNORMAL  = 1
	STATUS_USER_FORBIDDEN = 2
	// banner状态
	STATUS_BANNER_UP_LINE   = 1
	STATUS_BANNER_DOWN_LINE = 0
)

const (
	// banner类别
	TYPE_BANNER_HOME = 1
)

func init() {
	var err error
	MongoSession, err = mgo.Dial(global.MONGO_URI)
	if err != nil {
		panic(fmt.Sprintf("Error: connect to mongo db failed, error: %s", err))
	}
}

func GetMgoSession() *mgo.Session {
	return MongoSession.Copy()
}

/**
* 根据条件获取一个对象
 */
func GetModelFromMongo(dbName, collection string, query bson.M, model interface{}) (interface{}, error) {
	session := GetMgoSession()
	defer session.Close()

	modelInterface := GetReflectModel(model)
	if err := session.DB(dbName).C(collection).Find(query).One(modelInterface); err != nil {
		return nil, err
	}
	return GetReflectModel(modelInterface), nil
}

/**
* 根据条件获取多个对象
 */
func GetModelsFromMongo(dbName, collection string, query bson.M, arrayModel interface{}) (interface{}, error) {
	session := GetMgoSession()
	defer session.Close()

	arrayModelInterface := GetReflectModel(arrayModel)
	if err := session.DB(dbName).C(collection).Find(query).All(arrayModelInterface); err != nil {
		return nil, err
	}
	return GetReflectModel(arrayModelInterface), nil
}

/**
* 插入对象
 */
func InsertModelsToMongo(dbName, collection string, models ...interface{}) error {
	session := GetMgoSession()
	defer session.Close()

	if err := session.DB(dbName).C(collection).Insert(models...); err != nil {
		return err
	}
	return nil
}

/**
* 更新对象
 */
func UpdateModelFromMongo(dbName, collection string, selector, update bson.M) error {
	session := GetMgoSession()
	defer session.Close()

	if err := session.DB(dbName).C(collection).Update(selector, update); err != nil {
		return err
	}
	return nil
}

/**
* 更新插入对象
 */
func UpsertModelFromMongo(dbName, collection string, selector, update bson.M) (*mgo.ChangeInfo, error) {
	session := GetMgoSession()
	defer session.Close()

	info, err := session.DB(dbName).C(collection).Upsert(selector, update)
	if err != nil {
		return info, err
	}
	return info, nil
}

func GetReflectModel(model interface{}) interface{} {
	modelValue := reflect.ValueOf(model)
	return modelValue.Interface()
}
