package controller

import (
	"snack/controller/article"
	"snack/controller/index"
	"snack/controller/message"
	"snack/controller/user"
	middleware "snack/middleware/user"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

func GetRouter() *gin.Engine {
	return router
}

func init() {
	// 用户
	userController := user.UserController{}
	userApi := router.Group("/user")
	{
		// 用户登陆
		userApi.POST("/login", userController.UserLogin)
		// 用户注册
		userApi.POST("/regist", userController.UserRegist)

		// 获取用户信息
		userApi.GET("/info", middleware.JwtAuth(), userController.GetUserInfo)

		// 用户注销登陆
		userApi.PUT("/logout/:id", middleware.JwtAuth(), userController.UserLogout)

		// 用户列表
		userApi.GET("/list", middleware.JwtAuth(), userController.GetUserListByPage)
	}

	// 消息
	messageController := message.MessageController{}
	messageApi := router.Group("/message")
	{
		// 获取消息列表
		messageApi.GET("/list", middleware.JwtAuth(), messageController.GetMessageList)
		// 获取消息总数
		messageApi.GET("/count", middleware.JwtAuth(), messageController.GetMessageCount)
		// 标记消息已读
		messageApi.PUT("/read/:id", middleware.JwtAuth(), messageController.SetMessageRead)
	}

	// 主页
	indexController := index.IndexController{}
	indexApi := router.Group("/common")
	{
		// 获取主页Banner
		indexApi.GET("/banner/list", indexController.GetBanner)
	}

	// 文章
	articleController := article.ArticleController{}
	articleApi := router.Group("/article")
	{
		// 获取文章列表
		articleApi.GET("/list", articleController.GetArticleList)
		// 创建文章
		articleApi.POST("/create", articleController.CreateArticle)
	}
}
