package controller

import (
	"net/http"
	"temple-shrine-blog/model"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	blogs := model.GetAll()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"blogs": blogs,
	})
}

func ShowBlog(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetByID(id)
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"blog": blog,
	})
}
