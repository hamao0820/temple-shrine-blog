package controller

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"temple-shrine-blog/model"
	"temple-shrine-blog/util"

	"github.com/gin-gonic/gin"
)

const perPage = 6

func Index(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	var page int
	var err error
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit := perPage + 1
	offset := (page - 1) * perPage
	blogs := model.GetRange(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	count := len(blogs)
	sliceNum := perPage
	if count < sliceNum {
		sliceNum = count
	}
	blogs = blogs[:sliceNum]

	paginate := util.CreatePaginate(page, limit, count)

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
		"paginate":   paginate,
		"authorized": Authorized(c),
	})
}

func Awake(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
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
		"lat":       blog.Lat,
		"lng":       blog.Lng,
		"address":   blog.Address,
	}
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"blog":       showBlog,
		"authorized": Authorized(c),
		"api_key":    os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"authorized": Authorized(c),
	})
}

func ShowCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{
		"api_key": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

func ShowEdit(c *gin.Context) {
	id := c.Param("id")
	blog := model.GetOne(id)
	showBlog := map[string]any{
		"name":    blog.Name,
		"body":    blog.Body,
		"images":  blog.ImageURLs,
		"id":      blog.ID,
		"address": blog.Address,
		"lat":     blog.Lat,
		"lng":     blog.Lng,
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"blog":    showBlog,
		"api_key": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

func ShowNeedToLogin(c *gin.Context) {
	c.HTML(http.StatusUnauthorized, "need_to_login.html", gin.H{})
}

func Create(c *gin.Context) {
	blog := &model.Blog{
		Name:    c.PostForm("name"),
		Body:    c.PostForm("body"),
		Address: c.PostForm("address"),
	}
	lat, err := strconv.ParseFloat(c.PostForm("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid lat"})
		return
	}
	blog.Lat = lat
	lng, err := strconv.ParseFloat(c.PostForm("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid lng"})
		return
	}
	blog.Lng = lng

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
	blog.Address = c.PostForm("address")
	lat, err := strconv.ParseFloat(c.PostForm("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid lat"})
		return
	}
	blog.Lat = lat
	lng, err := strconv.ParseFloat(c.PostForm("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid lng"})
		return
	}
	blog.Lng = lng
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

	url, err := util.SaveImage(src)
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
