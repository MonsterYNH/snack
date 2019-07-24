package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ApiCount(c *gin.Context) {
	fmt.Println(c.Request.URL.Path)
	c.Next()
}
