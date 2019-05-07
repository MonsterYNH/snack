package main

import (
	"log"
	"net/http"
	"snack/controller"
	"snack/global"
)

func main() {
	// playerMap := playermap.CreatePlayerMap(20, 50)
	// playerMap.AddPlayer(player.CreatePlayer("test", node.Data{
	// 	X: 1,
	// 	Y: 1,
	// }))
	// playermap.PrintMap()
	// input := bufio.NewScanner(os.Stdin) //初始化一个扫表对象
	// for input.Scan() {
	// 	go func() {
	// 		if err := recover(); err != nil {
	// 			fmt.Println("Error: ", err)
	// 		}
	// 	}() //扫描输入内容
	// 	line := input.Text() //把输入内容转换为字符串
	// 	var operator int
	// 	flag := true
	// 	switch line {
	// 	case "w":
	// 		operator = player.PLAYER_DOWN
	// 	case "s":
	// 		operator = player.PLAYER_UP
	// 	case "a":
	// 		operator = player.PLAYER_RIGHT
	// 	case "d":
	// 		operator = player.PLAYER_LEFT
	// 	default:
	// 		fmt.Println("please input")
	// 		flag = false
	// 	}
	// 	if flag {
	// 		if err := playerMap.Play(playermap.Operat{
	// 			Operator:   operator,
	// 			PlayerName: "test",
	// 		}); err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 		playermap.GetMapInfo()
	// 		playermap.PrintMap()
	// 		//fmt.Println(playerMap.GetClientInfo("test"))
	// 	}

	// }
	log.Println("Server start on port ", global.SERVER_PORT)
	router := controller.GetRouter()
	http.ListenAndServe("0.0.0.0:"+global.SERVER_PORT, router)
}
