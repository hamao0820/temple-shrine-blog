package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(SECRET_KEY))
	r.Use(sessions.Sessions("LoginSession", store))
	r.LoadHTMLGlob("view/*.html")
	r.Static("images", "./images")
	r.GET("/", Index)
	r.GET("/blog/:id", ShowBlog)
	r.GET("/login", ShowLogin)
	r.POST("/login", Login)
	return r
}
