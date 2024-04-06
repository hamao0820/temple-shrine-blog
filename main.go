package main

import "temple-shrine-blog/controller"

func main() {
	r := controller.GetRouter()
	r.Run(":8080")
}
