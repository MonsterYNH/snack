package controller

import (
	"net/http"
	"snack/node"
	"snack/player"
	"snack/playermap"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start(url string) {
	router.POST("/snack/info", GetClientInfo)
	router.POST("/snack/play", Play)
	router.POST("/snack/addPlayer", AddPlayer)
	http.ListenAndServe(url, router)
}

type Parameter struct {
	Name   string `json:"name"`
	Operat int    `json:"operat"`
}

func GetClientInfo(c *gin.Context) {
	var parameter Parameter
	if err := c.ShouldBindJSON(&parameter); err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "parameter err"})
		return
	}
	data, err := playermap.GetClient(parameter.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func AddPlayer(c *gin.Context) {
	var parameter Parameter
	if err := c.ShouldBindJSON(&parameter); err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "parameter err"})
		return
	}
	player := player.CreatePlayer(parameter.Name, node.Data{
		X: 2,
		Y: 1,
	})
	playermap.GetMap().AddPlayer(player)
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
}

func Play(c *gin.Context) {
	var parameter Parameter
	if err := c.ShouldBindJSON(&parameter); err != nil {
		c.JSON(http.StatusOK, gin.H{"errmsg": "parameter err"})
		return
	}
	err := playermap.Play(parameter.Operat, parameter.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
}
