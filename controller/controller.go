package controller

import (
	"fmt"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"temple-shrine-blog/model"
	"temple-shrine-blog/util"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	blogs := model.GetAll()
	var showBlogs []map[string]any
	for _, b := range blogs {
		thumbnail := "images/thumbnail.png"
		if len(b.ImageURLs) > 0 {
			thumbnail = b.ImageURLs[0].URL
		}
		showblog := map[string]any{
			"thumbnail": thumbnail,
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
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	files := form.File["image_data"]
	for _, file := range files {
		url, err := saveImage(file, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot save image", "error": err})
			return
		}

		blog.ImageURLs = append(blog.ImageURLs, model.ImageURL{URL: url})
	}

	blog.Create()

	c.Redirect(http.StatusFound, "/")
}

func Edit(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	deleteImages := c.PostFormArray("delete-images[]")
	imageURLs := []model.ImageURL{}
	for _, i := range blog.ImageURLs {
		has := true
		for _, d := range deleteImages {
			if i.URL == d {
				has = false
				break
			}
		}
		if has {
			imageURLs = append(imageURLs, i)
		}
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}
	files := form.File["image_data"]
	for _, file := range files {
		url, err := saveImage(file, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot save image", "error": err})
			return
		}

		imageURLs = append(imageURLs, model.ImageURL{URL: url})
	}
	blog.Name = c.PostForm("name")
	blog.Body = c.PostForm("body")
	blog.ImageURLs = imageURLs
	blog.Edit()
	c.Redirect(http.StatusFound, fmt.Sprintf("/blog/%s", id))
}

func saveImage(file *multipart.FileHeader, c *gin.Context) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}

	url, err := util.SaveImage(img)
	if err != nil {
		return "", err
	}

	return url, nil
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	blog.Delete()
	c.Redirect(http.StatusFound, "/")
}
