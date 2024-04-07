package controller

import (
	"net/http"
	"temple-shrine-blog/model"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	blogs := model.GetAll()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"blogs":      blogs,
		"authorized": Authorized(c),
	})
}

func ShowBlog(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetByID(id)
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"blog":       blog,
		"authorized": Authorized(c),
	})
}

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"authorized": Authorized(c),
	})
}

func ShowCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{})
}

func ShowNeedToLogin(c *gin.Context) {
	c.HTML(http.StatusUnauthorized, "need_to_login.html", gin.H{})
}
