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
	var showBlogs []map[string]any
	for _, b := range blogs {
		showblog := map[string]any{
			"thumbnail": b.ImageURLs[0].URL,
			"createdAt": b.CreatedAt.Format("2006-01-02 15:04:05"),
			"name":      b.Name,
			"id":        b.ID,
		}
		showBlogs = append(showBlogs, showblog)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"blogs":      showBlogs,
		"authorized": Authorized(c),
	})
}

func ShowBlog(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	showBlog := map[string]any{
		"createdAt": blog.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt": blog.UpdatedAt.Format("2006-01-02 15:04:05"),
		"name":      blog.Name,
		"body":      blog.Body,
		"images":    blog.ImageURLs,
		"id":        blog.ID,
	}
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"blog":       showBlog,
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

func ShowEdit(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	showBlog := map[string]any{
		"name":   blog.Name,
		"body":   blog.Body,
		"images": blog.ImageURLs,
		"id":     blog.ID,
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"blog": showBlog,
	})
}

func ShowNeedToLogin(c *gin.Context) {
	c.HTML(http.StatusUnauthorized, "need_to_login.html", gin.H{})
}

func Create(c *gin.Context) {
	blog := &model.Blog{
		Name: c.PostForm("name"),
		Body: c.PostForm("body"),
	}
	fmt.Println(blog)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	files := form.File["image_data"]
	for _, file := range files {
		id, err := uuid.NewRandom()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		fileName := fmt.Sprintf("%s%s", id, filepath.Ext(file.Filename))
		p := filepath.Join("images", fileName)
		blog.ImageURLs = append(blog.ImageURLs, model.ImageURL{URL: p})

		c.SaveUploadedFile(file, p)
	}

	blog.Create()

	c.Redirect(http.StatusFound, "/")
}

func Edit(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	blog.Name = c.PostForm("name")
	blog.Body = c.PostForm("body")
	blog.Edit()
	c.Redirect(http.StatusFound, fmt.Sprintf("/blog/%s", id))
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	blog.Delete()
	c.Redirect(http.StatusFound, "/")
}
