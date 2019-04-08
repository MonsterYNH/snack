package article

import (
	"net/http"
	"snack/controller/common"
	Article "snack/model/article"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct{}

func (controller *ArticleController) GetArticleList(c *gin.Context) {
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
	list, err := Article.GetArticleListByTag(tag, start, limit)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(list, nil))
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
	if err := Article.CreateArticle(articleJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(true, nil))
}
