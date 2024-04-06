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
