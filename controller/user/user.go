package user

import (
	"encoding/json"
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
	id := c.Param("id")
	user, err := User.GetUser(bson.M{"_id": bson.ObjectIdHex(id), "status": db.STATUS_USER_NORMAL})

	var userInfoEntry UserInfoEntry
	bytes, _ := json.Marshal(user)
	if err := json.Unmarshal(bytes, &userInfoEntry); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.MARSHAML_ERROR))
		return
	}
	userInfoEntry.MessageCount, err = Message.GetUserMessageCount(userInfoEntry.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, common.ResponseSuccess(userInfoEntry, nil))
}

func (controller *UserController) UserLogin(c *gin.Context) {
	var loginEntry UserLoginEntry
	if err := c.BindJSON(&loginEntry); err != nil {
		log.Println("Error: user login failed, error: ", err)
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	query := make(map[string]interface{})
	query["password"] = loginEntry.Password
	switch loginEntry.Type {
	case "account":
		// TODO 账号校验 验证码校验
		query["account"] = loginEntry.Account
	case "email":
		// TODO 邮件校验 验证码校验
		query["email"] = loginEntry.Email
	case "phone":
		// TODO 手机校验 验证码校验
		query["phone"] = loginEntry.Phone
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
	bytes, _ := json.Marshal(user)
	var userInfo UserInfoEntry
	if err := json.Unmarshal(bytes, &userInfo); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.MARSHAML_ERROR))
		return
	}
	token, _ := middleware.CreateToken(middleware.CustomClaims{
		ID: userInfo.ID.Hex(),
	})

	c.JSON(http.StatusOK, common.ResponseSuccess(userInfo, token))
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
	if c.BindJSON(&userJson) != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.PARAMETER_ERR))
		return
	}
	if !userJson.Agreement {
		c.JSON(http.StatusOK, common.ResponseError(common.PLEASE_AGREEMENT))
		return
	}
	if userJson.Password != userJson.Confirm {
		c.JSON(http.StatusOK, common.ResponseError(common.DIFFERENT_PASSWORD))
		return
	}
	bytes, _ := json.Marshal(userJson)
	user := &User.User{}
	if err := json.Unmarshal(bytes, user); err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.MARSHAML_ERROR))
		return
	}
	var err error
	user, err = User.GetUser(bson.M{"phone": user.Phone, "status": 0})
	if err != nil && err != mgo.ErrNotFound {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	if user != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.USER_EXIST))
		return
	}
	user, err = User.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
		return
	}
	token, _ := middleware.CreateToken(middleware.CustomClaims{
		ID: user.ID.Hex(),
	})
	c.JSON(http.StatusOK, common.ResponseSuccess(user, token))
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
