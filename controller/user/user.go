package user

import (
	"log"
	"net/http"
	"snack/controller/common"
	"snack/db"
	middleware "snack/middleware/user"
	Message "snack/model/message"
	User "snack/model/user"
	"strconv"

	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct{}

func (controller *UserController) GetUserInfo(c *gin.Context) {
	if user := common.GetUser(c); user != nil {
		c.JSON(http.StatusOK, common.ResponseSuccess(user, nil))
		return
	} else {
		c.JSON(http.StatusOK, common.ResponseError(common.PLEASE_LOGIN))
		return
	}
}

func (controller *UserController) GetUserInfomation(c *gin.Context) {
	id := c.Param("id")
	user, err := User.GetUser(bson.M{"_id": bson.ObjectIdHex(id), "status": db.STATUS_USER_NORMAL})
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	// 获取关注信息
	if userId, exist := c.Get("user_id"); exist {
		if err := user.GetUserFollowStatus(userId.(string)); err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			return
		}
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(user, nil))

}

func (controller *UserController) UserLogin(c *gin.Context) {
	var loginEntry UserLoginEntry
	if err := c.ShouldBindJSON(&loginEntry); err != nil {
		log.Println("Error: user login failed, error: ", err)
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	query := make(map[string]interface{})
	switch loginEntry.Type {
	case "account":
		// TODO 账号校验 验证码校验
		query["account"] = loginEntry.Account
	case "email":
		// TODO 邮件校验 验证码校验
		query["email"] = loginEntry.Account
	case "phone":
		// TODO 手机校验 验证码校验
		query["phone"] = loginEntry.Account
	default:
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	user, err := User.GetUser(query)
	if err != nil {
		if err == mgo.ErrNotFound {
			c.JSON(http.StatusOK, common.ResponseError(common.USER_NOT_EXIST))
			return
		}
		log.Println("Error: user login failed, error: ", err, "parameter: ", loginEntry)
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	if user.Password != loginEntry.Password {
		c.JSON(http.StatusOK, common.ResponseError(common.PASSWORD_WRONG))
		return
	}
	token, _ := middleware.CreateToken(middleware.CustomClaims{
		ID: user.ID.Hex(),
	})
	c.JSON(http.StatusOK, common.ResponseSuccess(user, token))
}

func (controller *UserController) UserLogout(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		c.JSON(http.StatusOK, common.ResponseError(common.ID_NOT_EXIST))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess("success", nil))
}

func (controller *UserController) UserRegist(c *gin.Context) {
	var userJson UserRegistEntry
	if c.ShouldBindJSON(&userJson) != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	newUser := &User.User{}

	if userJson.Password != userJson.Confirm && len(userJson.Password) == 0 {
		c.JSON(http.StatusOK, common.ResponseError(common.DIFFERENT_PASSWORD))
		return
	}
	query := bson.M{
		"status": 0,
	}
	switch userJson.Type {
	case "account":
		query["account"] = userJson.Account
		newUser.Account = userJson.Account
	case "phone":
		query["phone"] = userJson.Account
		newUser.Phone = userJson.Account
	case "email":
		query["email"] = userJson.Account
		newUser.Email = userJson.Account
	}
	newUser.Password = userJson.Password
	user, err := User.GetUser(query)
	if err != nil {
		if err != mgo.ErrNotFound {
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			return
		}
	}
	if bson.IsObjectIdHex(user.ID.Hex()) {
		c.JSON(http.StatusOK, common.ResponseError(common.USER_EXIST))
		return
	}
	newUser, err = User.CreateUser(*newUser)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	if _, err = Message.SaveMessage(newUser.ID, bson.NewObjectId(), "欢迎加入我们", Message.SYSTEM_MESSAGE, nil); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	token, err := middleware.CreateToken(middleware.CustomClaims{
		ID: newUser.ID.Hex(),
	})
	c.JSON(http.StatusOK, common.ResponseSuccess(newUser, token))
	return
}

func (controller *UserController) GetUserListByPage(c *gin.Context) {
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
	users, err := User.GetUsers(nil, start, limit)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	total, err := User.GetUserCount(nil)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	userPage := UserPage{
		Total: total,
		Data:  users,
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(userPage, nil))
	return
}

func (controller *UserController) FollowUser(c *gin.Context) {
	var userFollowJson UserFollowEntry
	if err := c.ShouldBindJSON(&userFollowJson); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	operatorId, _ := c.Get("user_id")
	userId := c.Param("id")

	if err := User.FollowUser(userFollowJson.Option, bson.ObjectIdHex(userId), bson.ObjectIdHex(operatorId.(string))); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(gin.H{"status": true}, nil))
}
