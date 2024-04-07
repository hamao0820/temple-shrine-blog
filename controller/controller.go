package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"temple-shrine-blog/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func Create(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}
	name := c.PostForm("name")
	fmt.Println(name)
	description := c.PostForm("description")
	fmt.Println(description)

	files := form.File["image_data"]
	for _, file := range files {
		id, err := uuid.NewRandom()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		p := filepath.Join("images", fmt.Sprintf("%s%s", id, filepath.Ext(file.Filename)))

		c.SaveUploadedFile(file, p)
	}

	c.Redirect(http.StatusFound, "/")
}
