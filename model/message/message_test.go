package model

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestMessage(t *testing.T) {
	GetUserMessageCount(bson.ObjectIdHex("5c98aeaafb7d8c4e73a1f712"))
}

func TestCreateMessage(t *testing.T) {
	objectId := bson.ObjectIdHex("5c98aeaafb7d8c4e73a1f712")
	message, err := CreateUserMessage(bson.ObjectIdHex("5c98aeaafb7d8c4e73a1f712"), bson.ObjectIdHex("5c98aeaafb7d8c4e73a1f712"), "this is a test", 1, &objectId)
	if err != nil {
		t.Error(err)
	}
	log.Print(message)
}
