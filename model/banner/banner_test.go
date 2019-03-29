package model

import (
	"snack/db"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestBanner(t *testing.T) {
	id := bson.NewObjectId()
	banner, err := CreateBanner(db.TYPE_BANNER_HOME, &id, "www.spider.store")
	t.Logf("%v, %v", banner, err)
	if err != nil {
		t.Error(err)
	}
}
