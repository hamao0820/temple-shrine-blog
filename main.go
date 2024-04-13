package main

import (
	"log"
	"os"
	"temple-shrine-blog/controller"
)

func main() {
	r := controller.GetRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
