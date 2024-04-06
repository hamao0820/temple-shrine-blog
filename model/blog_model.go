package model

import (
	"time"
)

type Blog struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	Thumbnail string   `json:"thumbnail"`
	CreatedAt string   `json:"created_at"`
}

func GetAll() []Blog {
	blogs := []Blog{
		{
			ID:        "1",
			Title:     "Title 1",
			Content:   "Content 1",
			Images:    []string{"images/sample1.png", "images/sample2.png"},
			Thumbnail: "images/sample1.png",
			CreatedAt: time.Now().Format("2006-01-02"),
		},
		{
			ID:        "2",
			Title:     "Title 2",
			Content:   "Content 2",
			Images:    []string{"images/sample3.png", "images/sample4.png"},
			Thumbnail: "images/sample3.png",
			CreatedAt: time.Now().Format("2006-01-02"),
		},
	}
	return blogs
}
