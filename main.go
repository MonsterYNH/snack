package main

import (
	"log"
	"net/http"
	"snack/controller"
	"snack/global"
)

func main() {
	log.Println("Server start on port ", global.SERVER_PORT)
	router := controller.GetRouter()
	http.ListenAndServe("0.0.0.0:"+global.SERVER_PORT, router)
}
