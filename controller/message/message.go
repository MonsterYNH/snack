package message

import (
	"net/http"
	"snack/controller/common"
	Message "snack/model/message"
	User "snack/model/user"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

type MessageController struct{}

func (message *MessageController) GetMessageList(c *gin.Context) {
	userInfo, _ := c.Get("user")
	userInfoEntry, _ := userInfo.(*User.User)
	start, err := strconv.Atoi(c.DefaultQuery("start", "1"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	messages, err := Message.GetUserMessage(userInfoEntry.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	total, err := Message.GetUserMessageCount(userInfoEntry.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	pageEntry := MessagePageEntry{
		Total:   total,
		HasNext: start*limit < total,
		Data:    messages,
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(pageEntry, nil))
}

func (message *MessageController) GetMessageCount(c *gin.Context) {
	user, _ := c.Get("user")
	userEntry, exist := user.(User.User)
	if !exist {
		c.JSON(http.StatusOK, common.ResponseError(common.USER_NOT_EXIST))
		return
	}
	count, err := Message.GetUserMessageCount(userEntry.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{"count": count}, nil))
}

func (message *MessageController) SetMessageRead(c *gin.Context) {
	id := c.Param("id")
	err := Message.SetUserMessageRead(bson.ObjectIdHex(id))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(true, nil))
}
