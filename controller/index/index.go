package index

import (
	"net/http"
	"snack/controller/common"
	"snack/db"
	model "snack/model/banner"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (index *IndexController) GetBanner(c *gin.Context) {
	banners, err := model.GetBannerByType(db.TYPE_BANNER_HOME)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(banners, nil))
	return
}
