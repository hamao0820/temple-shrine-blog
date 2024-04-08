package model

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const DBName = "test.db"

type Blog struct {
	gorm.Model
	Name      string
	Body      string
	Address   string
	Lat       float64
	Lng       float64
	ImageURLs []ImageURL `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,foreignkey:OrganizationID;"`
}

type ImageURL struct {
	gorm.Model
	URL    string
	BlogID uint
}

func init() {
	db := openDB()

	db.AutoMigrate(&Blog{}, &ImageURL{})
}

func openDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", DBName)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func GetAll() (datas []Blog) {
	db := openDB()
	res := db.Model(&Blog{}).Preload("ImageURLs").Find(&datas).Order("created_at desc")
	if res.Error != nil {
		return
	}
	return datas
}

func GetOne(id string) (data Blog) {
	db := openDB()
	result := db.Model(&Blog{}).Preload("ImageURLs").First(&data, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return data
}

func (b *Blog) Create() {
	db := openDB()
	r := db.Create(b)
	if r.Error != nil {
		panic(r.Error)
	}
}

func (b *Blog) Edit() {
	db := openDB()
	var urls []ImageURL
	if r := db.Where("blog_id = ?", b.ID).Find(&urls); r.Error != nil {
		panic(r.Error)
	}
	if r := db.Save(b); r.Error != nil {
		panic(r.Error)
	}
	for i := range urls {
		url := urls[i]
		has := false
		for _, i := range b.ImageURLs {
			if url.ID == i.ID {
				has = true
				break
			}
		}
		if !has {
			deleteImage(url.URL)
			if r := db.Delete(&url); r.Error != nil {
				panic(r.Error)
			}
		}
	}
}

func (b *Blog) Delete() {
	db := openDB()
	if r := db.Delete(b); r.Error != nil {
		panic(r.Error)
	}
	for _, i := range b.ImageURLs {
		deleteImage(i.URL)
	}
	if r := db.Where("blog_id = ?", b.ID).Delete(&ImageURL{}); r.Error != nil {
		panic(r.Error)
	}
}

func deleteImage(name string) {
	os.Remove(name)
}
