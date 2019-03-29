package controller

import (
	"snack/controller/index"
	"snack/controller/message"
	"snack/controller/user"
	middleware "snack/middleware/user"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.Use(middleware.JwtAuth())
}

func GetRouter() *gin.Engine {
	return router
}

func init() {
	userController := user.UserController{}
	userApi := router.Group("/user")
	{
		// 用户登陆
		userApi.POST("/login", userController.UserLogin)

		// 获取用户信息
		userApi.GET("/info/:id", userController.GetUserInfo)

		// 用户注销登陆
		userApi.PUT("/logout/:id", userController.UserLogout)

		// 用户注册
		userApi.POST("/regist", userController.UserRegist)

		// 用户列表
		userApi.GET("/list", userController.GetUserListByPage)
	}

	messageController := message.MessageController{}
	messageApi := router.Group("/message")
	{
		// 获取消息列表
		messageApi.GET("/list/:id", messageController.GetMessageList)
		// 获取消息总数
		messageApi.GET("/count", messageController.GetMessageCount)
	}

	indexController := index.IndexController{}
	indexApi := router.Group("/common")
	{
		// 获取主页Banner
		indexApi.GET("/banner/list", indexController.GetBanner)
	}
}
