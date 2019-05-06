package article

import (
	"errors"
	"net/http"
	"snack/controller/common"
	Article "snack/model/article"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

type ArticleController struct{}

// func (controller *ArticleController) GetArticleList(c *gin.Context) {
// 	var err error
// 	tag := c.Query("tag")
// 	start, err := strconv.Atoi(c.DefaultQuery("start", "1"))
// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
// 		return
// 	}
// 	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
// 		return
// 	}

// 	// 验证是否有用户
// 	result := make(map[string]interface{})
// 	userIDStr, exist := c.Get("user_id")
// 	if exist {
// 		userID := bson.ObjectIdHex(userIDStr.(string))
// 		result, err = Article.GetArticleListByTag(tag, &userID, start, limit)
// 	} else {
// 		result, err = Article.GetArticleListByTag(tag, nil, start, limit)
// 	}

// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, common.ResponseSuccess(result, nil))
// }

func (controller *ArticleController) GetArticleById(c *gin.Context) {
	id := c.Param("id")
	userIDStr, exist := c.Get("user_id")
	var (
		article *Article.ArticleEntry
		err     error
	)

	if exist {
		userID := bson.ObjectIdHex(userIDStr.(string))
		article, err = Article.GetArticleById(bson.ObjectIdHex(id), &userID)
	} else {
		article, err = Article.GetArticleById(bson.ObjectIdHex(id), nil)
	}

	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(article, nil))
}

func (controllere *ArticleController) CreateArticle(c *gin.Context) {
	var articleJson Article.Article
	if err := c.ShouldBindJSON(&articleJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	if len(articleJson.Title) == 0 {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, errors.New("article title can not be empty")))
		return
	}
	if len(articleJson.Content) == 0 {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, errors.New("article content can not be empty")))
		return
	}
	if userId, exist := c.Get("user_id"); exist {
		articleJson.CreateUser = bson.ObjectIdHex(userId.(string))
	}
	if err := Article.CreateArticle(articleJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(true, nil))
}

func (controller *ArticleController) GetArticleTags(c *gin.Context) {
	tags, err := Article.GetArticleTags()
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(tags, nil))
}

func (controller *ArticleController) GetArticleByType(c *gin.Context) {
	var err error
	articleType := c.Query("type")
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
	sort, err := strconv.Atoi(c.DefaultQuery("sort", "-1"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	var (
		userID *bson.ObjectId
	)
	userIDStr, exist := c.Get("user_id")
	if exist {
		userIDTemp := bson.ObjectIdHex(userIDStr.(string))
		userID = &userIDTemp
	} else {
		userID = nil
	}

	result := make(map[string]interface{}, 0)
	switch articleType {
	case "new":
		result, err = Article.GetNewArticle(userID, start, limit, sort)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	case "hot":
		result, err = Article.GetHotArticle(userID, start, limit, sort)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	case "recommend":
		result, err = Article.GetRecommendArticle(userID, start, limit)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	case "all":
		tags := c.QueryArray("tags[]")
		result, err = Article.GetArticleListByTag(tags, nil, start, limit, sort)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	case "same":
		id := c.Query("id")
		result, err = Article.GetSameArticle(bson.ObjectIdHex(id), start, limit)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	case "user":
		id := c.Query("user_id")
		result, err = Article.GetArticleByUserId(bson.ObjectIdHex(id), start, limit)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
			return
		}
	default:
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(result, nil))
}

func (controller *ArticleController) GetArticleCategories(c *gin.Context) {
	categories, err := Article.GetArticleCategories()
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(categories, nil))
}

// func (controller *ArticleController) GetNewArticleComment(c *gin.Context) {
// 	start, err := strconv.Atoi(c.DefaultQuery("start", "1"))
// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
// 		return
// 	}
// 	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
// 		return
// 	}
// 	result, err := Article.GetNewArticleComment(start, limit)
// 	if err != nil {
// 		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR, err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, common.ResponseSuccess(result, nil))
// }
