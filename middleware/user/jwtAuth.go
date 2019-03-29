package middleware

import (
	"fmt"
	"net/http"
	"snack/controller/common"
	"snack/db"
	model "snack/model/user"
	"strings"
	"time"

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
		}
		if strings.HasPrefix(c.Request.RequestURI, "/user/login") && c.Request.Method == "POST" {
			c.Next()
		}
		token := c.Request.Header.Get("authorization")
		customClaims, code := ParseToken(token)
		if code > 0 {
			c.JSON(http.StatusOK, common.ResponseError(code))
			c.Abort()
		}
		if !bson.IsObjectIdHex(customClaims.ID) {
			c.JSON(http.StatusOK, common.ResponseError(common.ID_NOT_EXIST))
			c.Abort()
		}
		user, err := model.GetUser(bson.M{"_id": bson.ObjectIdHex(customClaims.ID), "status": db.STATUS_USER_NORMAL})
		if err != nil {
			c.JSON(http.StatusOK, common.ResponseError(common.USER_NOT_EXIST))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(123)
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
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return claims, common.SUCCESS
	}
	return nil, common.TOKEN_INVALID
}
