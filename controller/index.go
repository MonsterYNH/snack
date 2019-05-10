package controller

import (
	"snack/controller/article"
	"snack/controller/community"
	"snack/controller/index"
	"snack/controller/leavemessage"
	"snack/controller/message"
	"snack/controller/user"
	"snack/controller/wechat"
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
	userApi := router.Group("/api/user")
	{
		// 用户登陆
		userApi.POST("/login", userController.UserLogin)
		// 用户注册
		userApi.POST("/regist", userController.UserRegist)

		// 获取用户信息
		userApi.GET("/info", middleware.JwtAuth(), userController.GetUserInfo)
		userApi.GET("/info/:id", middleware.WithUser(), userController.GetUserInfomation)

		// 用户注销登陆
		userApi.PUT("/logout/:id", middleware.JwtAuth(), userController.UserLogout)

		// 用户列表
		userApi.GET("/list", middleware.JwtAuth(), userController.GetUserListByPage)

		// 用户关注
		userApi.PUT("/follow/:id", middleware.JwtAuth(), userController.FollowUser)
		// 用户关注列表
		userApi.GET("/followed/:id", middleware.WithUser(), userController.GetUserFollowed)
	}
	wechatController := wechat.WxConnectController{}
	router.GET("/api/wx", wechatController.Get)

	// 消息
	messageController := message.MessageController{}
	messageApi := router.Group("/api/message")
	{
		// 获取消息列表
		messageApi.GET("/list", middleware.JwtAuth(), messageController.GetMessageList)
		// 获取消息总数
		messageApi.GET("/count", middleware.JwtAuth(), messageController.GetMessageCount)
		// 标记消息已读
		messageApi.PUT("/read/:id", middleware.JwtAuth(), messageController.SetMessageRead)
		// 标记消息全部已读
		messageApi.PUT("/readall", middleware.JwtAuth(), messageController.SetMessageReadAll)
	}

	// 主页
	indexController := index.IndexController{}
	indexApi := router.Group("/api/common")
	{
		// 获取主页Banner
		indexApi.GET("/banner/list", indexController.GetBanner)
	}

	// 文章
	articleController := article.ArticleController{}
	articleApi := router.Group("/api/article")
	{
		// 获取文章列表
		//articleApi.GET("/list", middleware.WithUser(), articleController.GetArticleList)
		// 创建文章
		articleApi.POST("/create", middleware.JwtAuth(), articleController.CreateArticle)
		// 获取文章标签
		articleApi.GET("/tags", articleController.GetArticleTags)
		// 获取文章类别
		articleApi.GET("/categories", articleController.GetArticleCategories)
		// 获取文章列表
		articleApi.GET("/type", middleware.WithUser(), articleController.GetArticleByType)
		// 获取文章详情
		articleApi.GET("/detail/:id", middleware.WithUser(), articleController.GetArticleById)
	}

	// 社区
	communityController := community.CommunityController{}
	communityApi := router.Group("/api/community")
	{
		// 社区通用操作
		communityApi.POST("/option/:id", middleware.JwtAuth(), communityController.CommunityOption)
		// 创建评论
		communityApi.POST("/create/comment", middleware.JwtAuth(), communityController.CreateComment)
		// 创建评论回复
		communityApi.POST("/create/commentReply", middleware.JwtAuth(), communityController.CreateCommentReply)
		// 获取评论列表
		communityApi.GET("/comment/list", middleware.WithUser(), communityController.GetCommentByPage)
		// 获取回复列表
		communityApi.GET("/reply/list", communityController.GetCommentRepliesByPage)
	}

	// 留言
	leaveMessageController := leavemessage.LeaveMessageController{}
	leaveMessageApi := router.Group("/api/leaveMessage")
	{
		// 获取留言列表
		leaveMessageApi.GET("/list", middleware.WithUser(), leaveMessageController.GetLeaveMessages)
		// 留言
		leaveMessageApi.POST("/create/leaveMessage", middleware.WithUser(), leaveMessageController.CreateLeaveMessage)
		// 评论留言
		leaveMessageApi.POST("/comment", middleware.WithUser(), leaveMessageController.CommentLeaveMessage)
	}
}
