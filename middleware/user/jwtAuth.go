package middleware

import (
	"net/http"
	"snack/controller/common"
	"snack/db"
	User "snack/model/user"
	"strings"

	mgo "gopkg.in/mgo.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

const jwtKey = "mt_jwt_test"

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/user/regist" && c.Request.Method == "POST" {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.RequestURI, "/user/login") && c.Request.Method == "POST" {
			c.Next()
			return
		}
		token := c.Request.Header.Get("authorization")
		customClaims, code := ParseToken(token)
		if code > 0 && customClaims == nil {
			c.JSON(http.StatusOK, common.ResponseError(code))
			c.Abort()
			return
		}
		if !bson.IsObjectIdHex((*customClaims).ID) {
			c.JSON(http.StatusOK, common.ResponseError(common.ID_NOT_EXIST))
			c.Abort()
			return
		}
		user, err := User.GetUser(bson.M{"_id": bson.ObjectIdHex(customClaims.ID), "status": db.STATUS_USER_NORMAL})
		if err != nil {
			if err == mgo.ErrNotFound {
				c.JSON(http.StatusOK, common.ResponseError(common.USER_NOT_EXIST))
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, common.ResponseError(common.SERVER_ERROR))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenString string) (*CustomClaims, int) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, common.TOKEN_INVALID
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, common.TOKEN_EXPIRED
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, common.TOKEN_INVALID
			} else {
				return nil, common.TOKEN_INVALID
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, common.SUCCESS
	}
	return nil, common.TOKEN_INVALID
}
