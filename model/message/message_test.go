package model

import (
	"encoding/json"
	"fmt"
	"snack/db"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestMessage(t *testing.T) {
	GetUserMessageCount(bson.ObjectIdHex("5c98aeaafb7d8c4e73a1f712"))
}

func TestCreateMessage(t *testing.T) {
	objectId := bson.ObjectIdHex("5ca0834bfb7d8c30a495b889")
	session := db.GetMgoSession()
	defer session.Close()

	messageInfos := make([]MessageInfo, 0)
	aggregate := []bson.M{
		bson.M{"$match": bson.M{"user_id": objectId}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "from_user_id", "foreignField": "_id", "as": "user_info"}},
		bson.M{"$unwind": "$user_info"},
		bson.M{"$project": bson.M{"from_user_id": "$from_user_id", "from_user_avatar": "$user_info.avatar", "object_id": "$user_info.object_id", "title": "$title", "content": "$content", "type": "$type"}},
	}
	bytes, _ := json.Marshal(aggregate)
	fmt.Println(string(bytes))
	if err := session.DB(db.DB_NAME).C(db.COLL_MESSAGE).Pipe(aggregate).All(&messageInfos); err != nil {
		t.Error(err)
	}
	bytes, _ = json.Marshal(messageInfos)
	fmt.Println(string(bytes))
}
