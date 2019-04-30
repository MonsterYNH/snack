package community

import (
	"errors"
	"fmt"
	"net/http"
	"snack/controller/common"
	Article "snack/model/article"
	Comment "snack/model/comment"
	Message "snack/model/message"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type CommunityController struct{}

func (controller *CommunityController) CommunityOption(c *gin.Context) {
	var optionJson CommunityOption
	if err := c.ShouldBindJSON(&optionJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	var (
		//err    error
		num    int
		status int
	)

	switch optionJson.Kind {
	case "article":
		article, err := Article.GetArticleById(bson.ObjectIdHex(id))
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
		switch optionJson.Type {
		case "like":
			if num, status, err = Article.LikeArticle(bson.ObjectIdHex(id), bson.ObjectIdHex(userID.(string)), optionJson.Option); err != nil {
				c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
				return
			}
			if _, err = Message.SaveMessage(article.CreateUser, bson.ObjectIdHex(userID.(string)), "去查看", Message.USER_ARTICLE_LIKE, nil); err != nil {
				c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
				return
			}
		case "collect":
			if num, status, err = Article.CollectArticle(bson.ObjectIdHex(id), bson.ObjectIdHex(userID.(string)), optionJson.Option); err != nil {
				c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
				return
			}
			if _, err = Message.SaveMessage(article.CreateUser, bson.ObjectIdHex(userID.(string)), "去查看", Message.USER_ARTICLE_COLLECT, nil); err != nil {
				c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
				return
			}
		default:
			c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, errors.New("no match community option")))
			return
		}
	case "comment":
		if err := Comment.LikeComment(bson.ObjectIdHex(id), bson.ObjectIdHex(userID.(string))); err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
		// if _, err = Message.SaveMessage(article.CreateUser, bson.ObjectIdHex(userID.(string)), "去查看", Message.USER_ARTICLE_COMMENT, nil); err != nil {
		// 	c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		// 	return
		// }
	case "commentReply":
		// if _, err = Message.SaveMessage(article.CreateUser, bson.ObjectIdHex(userID.(string)), "去查看", Message.USER_ARTICLE_COMMENT_REPLY, nil); err != nil {
		// 	c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		// 	return
		// }
	default:
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, errors.New("no match community option")))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{
		"num":    num,
		"status": status,
	}, nil))
}

func (controller *CommunityController) CreateComment(c *gin.Context) {
	var createJson CommunityCreate
	if err := c.ShouldBindJSON(&createJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	id, _ := c.Get("user_id")

	if err := Comment.CreateComment(createJson.Kind, createJson.Content, createJson.ObjectId, bson.ObjectIdHex(id.(string))); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	num, err := Comment.GetCommentNum(createJson.Kind, createJson.ObjectId)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{"comment_num": num}, nil))
}

func (controller *CommunityController) CreateCommentReply(c *gin.Context) {
	var createJson CommunityCreate
	if err := c.ShouldBindJSON(&createJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	id, _ := c.Get("user_id")

	if err := Comment.CreateCommentReply(createJson.ObjectId, createJson.CommentUser, bson.ObjectIdHex(id.(string)), createJson.Content); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(true, nil))
}

func (controller *CommunityController) GetCommentByPage(c *gin.Context) {
	kind := c.Query("kind")
	objectId := c.Query("object_id")
	var userId *bson.ObjectId
	if userIdStr, exist := c.Get("user_id"); exist {
		id := bson.ObjectIdHex(userIdStr.(string))
		userId = &id
	} else {
		userId = nil
	}

	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}

	list, err := Comment.GetCommentsByPage(kind, bson.ObjectIdHex(objectId), userId, start, limit, 5)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(list, nil))
}

func (controller *CommunityController) GetCommentRepliesByPage(c *gin.Context) {
	id := c.Query("id")
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}

	list, err := Comment.GetCommentRepliesByPage(bson.ObjectIdHex(id), start, limit)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(list, nil))
}
