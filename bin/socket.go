package main

import (
	"log"
	"net/http"
	"snack/controller"
)

func main() {
	log.Println("Server start on port 5000")
	router := controller.GetRouter()
	http.ListenAndServe("localhost:5000", router)

}
