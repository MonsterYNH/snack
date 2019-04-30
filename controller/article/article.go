package article

import (
	"net/http"
	"snack/controller/common"
	Article "snack/model/article"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

type ArticleController struct{}

func (controller *ArticleController) GetArticleList(c *gin.Context) {
	var err error
	tag := c.Query("tag")
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

	// 验证是否有用户
	list := make([]Article.ArticleEntry, 0)
	userIDStr, exist := c.Get("user_id")
	if exist {
		userID := bson.ObjectIdHex(userIDStr.(string))
		list, err = Article.GetArticleListByTag(tag, &userID, start, limit)
	} else {
		list, err = Article.GetArticleListByTag(tag, nil, start, limit)
	}

	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	total, err := Article.GetArticleCountByTag(tag)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(ArticleListEntry{
		Data:  list,
		Total: total,
	}, nil))
}

func (controllere *ArticleController) CreateArticle(c *gin.Context) {
	var articleJson Article.Article
	if err := c.ShouldBindJSON(&articleJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	if len(articleJson.Title) == 0 {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	if len(articleJson.Content) == 0 {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	if userId, exist := c.Get("user_id"); exist {
		articleJson.CreateUser = bson.ObjectIdHex(userId.(string))
	}
	if err := Article.CreateArticle(articleJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(true, nil))
}

func (controller *ArticleController) GetArticleTags(c *gin.Context) {
	tags, err := Article.GetArticleTags()
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(tags, nil))
}

func (controller *ArticleController) GetArticleByType(c *gin.Context) {
	var err error
	articleType := c.Query("type")
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
	sort, err := strconv.Atoi(c.DefaultQuery("sort", "-1"))
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}

	result := make(map[string]interface{}, 0)
	switch articleType {
	case "new":
		result, err = Article.GetNewArticle(start, limit, sort)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			return
		}
	case "hot":
		result, err = Article.GetHotArticle(start, limit, sort)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			return
		}
	case "recommend":
		result, err = Article.GetRecommendArticle(start, limit)
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			return
		}
	default:
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(result, nil))
}
