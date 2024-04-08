package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const DBName = "test.db"

type Blog struct {
	gorm.Model
	Name      string
	Body      string
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

func (b *Blog) Delete() {
	db := openDB()
	if r := db.Delete(b); r.Error != nil {
		panic(r.Error)
	}
	if r := db.Where("blog_id = ?", b.ID).Delete(&ImageURL{}); r.Error != nil {
		panic(r.Error)
	}
}
