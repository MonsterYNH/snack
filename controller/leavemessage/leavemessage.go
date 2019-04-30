package leavemessage

import (
	"log"
	"net/http"
	"snack/controller/common"
	LeaveMessage "snack/model/leavemessage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaveMessageController struct{}

func (controller *LeaveMessageController) GetLeaveMessages(c *gin.Context) {
	start, err := strconv.Atoi(c.DefaultQuery("start", "1"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	list, err := LeaveMessage.GetLeaveMessage(start, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(list, nil))
}

func (controller *LeaveMessageController) CreateLeaveMessage(c *gin.Context) {
	var createLessageMessageJson CreateLeaveMessageEntry
	if err := c.ShouldBindJSON(&createLessageMessageJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	if err := LeaveMessage.CreateLeaveMessage(createLessageMessageJson.Content, createLessageMessageJson.CreateUserName, createLessageMessageJson.CreateUserAvatar, createLessageMessageJson.CreateUser); err != nil {
		log.Println("ERROR create lessage message, error: ", err)
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{"data": true}, nil))
}

func (controller *LeaveMessageController) CommentLeaveMessage(c *gin.Context) {
	var commentLeaveMessageJson CommentLeaveMessageEntry
	if err := c.ShouldBindJSON(&commentLeaveMessageJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	if err := LeaveMessage.CommentLeaveMessage(commentLeaveMessageJson.MessageId, commentLeaveMessageJson.Content, commentLeaveMessageJson.CreateUserName, commentLeaveMessageJson.CreateUserAvatar, commentLeaveMessageJson.CreateUser); err != nil {
		log.Println("ERROR create lessage message, error: ", err)
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{"data": true}, nil))
}
