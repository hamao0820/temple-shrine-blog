package model

import (
	"fmt"
	"os"
	"temple-shrine-blog/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Blog struct {
	gorm.Model
	Name       string
	Body       string
	Address    string
	Lat        float64
	Lng        float64
	ImageNames []ImageName `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL,foreignkey:OrganizationID;"`
}

type ImageName struct {
	gorm.Model
	Name   string
	BlogID uint
}

func init() {
	db := openDB()

	db.AutoMigrate(&Blog{}, &ImageName{})
}

func openDB() *gorm.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	userName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, userName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func GetAll() (datas []Blog) {
	db := openDB()
	res := db.Model(&Blog{}).Preload("ImageNames").Find(&datas).Order("created_at desc")
	if res.Error != nil {
		return
	}
	return datas
}

func GetOne(id string) (data Blog) {
	db := openDB()
	result := db.Model(&Blog{}).Preload("ImageNames").First(&data, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return data
}

func GetRange(limit, offset int) (datas []Blog) {
	db := openDB()
	res := db.Model(&Blog{}).Preload("ImageNames").Limit(limit).Offset(offset).Order("created_at desc").Find(&datas)
	if res.Error != nil {
		return
	}
	return datas
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
	var names []ImageName
	if r := db.Where("blog_id = ?", b.ID).Find(&names); r.Error != nil {
		panic(r.Error)
	}
	if r := db.Save(b); r.Error != nil {
		panic(r.Error)
	}
	for i := range names {
		name := names[i]
		has := false
		for _, i := range b.ImageNames {
			if name.ID == i.ID {
				has = true
				break
			}
		}
		if !has {
			util.DeleteImage(name.Name)
			if r := db.Delete(&name); r.Error != nil {
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
	for _, i := range b.ImageNames {
		err := util.DeleteImage(i.Name)
		if err != nil {
			panic(err)
		}
	}
	if r := db.Where("blog_id = ?", b.ID).Delete(&ImageName{}); r.Error != nil {
		panic(r.Error)
	}
}
