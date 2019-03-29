package model

import (
	"snack/db"

	"gopkg.in/mgo.v2/bson"
)

type Banner struct {
	ID       bson.ObjectId  `json:"-" bson:"_id"`
	Type     int            `json:"-" bson:"type"`
	ObjectID *bson.ObjectId `json:"object_id" bson:"object_id"`
	Link     string         `json:"link" bson:"link"`
	Status   int            `json:"-" bson:"status"`
}

func CreateBanner(bannerType int, objectID *bson.ObjectId, link string) (*Banner, error) {
	session := db.GetMgoSession()
	defer session.Close()

	banner := &Banner{
		ID:       bson.NewObjectId(),
		Type:     bannerType,
		ObjectID: objectID,
		Link:     link,
		Status:   db.STATUS_BANNER_UP_LINE,
	}
	err := session.DB(db.DB_NAME).C(db.COLL_BANNER).Insert(banner)
	return banner, err
}

func GetBannerByType(bannerType int) ([]Banner, error) {
	session := db.GetMgoSession()
	defer session.Close()

	banners := make([]Banner, 0)
	err := session.DB(db.DB_NAME).C(db.COLL_BANNER).Find(bson.M{"type": bannerType, "status": db.STATUS_BANNER_UP_LINE}).All(&banners)
	return banners, err
}
