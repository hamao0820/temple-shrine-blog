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
	ImageURLs []ImageURL `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ImageURL struct {
	gorm.Model
	URL    string
	BlogID uint
}

func init() {
	db := openDB()

	db.AutoMigrate(&Blog{})
	db.AutoMigrate(&ImageURL{})
}

func openDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", DBName)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// func GetAll() (datas []Blog) {
// 	res := db.Find(&datas)
// 	if res.Error != nil {
// 		return
// 	}
// 	return datas
// }

// func GetOne(id int) (data Blog) {
// 	result := db.First(&data, id)
// 	if result.Error != nil {
// 		panic(result.Error)
// 	}
// 	return data
// }

func (b *Blog) Create() {
	db := openDB()
	r := db.Create(b)
	if r.Error != nil {
		panic(r.Error)
	}
}
