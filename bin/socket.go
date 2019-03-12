package main

import (
	"log"
	"snack/controller"
)

func main() {
	log.Println("Server start on port 5000")
	controller.Start("0.0.0.0:5000")
}
