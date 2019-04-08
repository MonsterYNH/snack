package common

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"snack/db"
	User "snack/model/user"
)

func GetUser(c *gin.Context) *User.User {
	userID, exist := c.Get("user_id")
	if !exist {
		return nil
	}
	user, err := User.GetUser(bson.M{"_id": bson.ObjectIdHex(userID.(string)), "status": db.STATUS_USER_NORMAL})
	if err != nil {
		return nil
	}
	return user
}
