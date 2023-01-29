package main

import (
	"simple_tiktok/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main9() {
	db, err := gorm.Open(mysql.Open("root:124578@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.VideoBasic{})
	db.AutoMigrate(&models.CommentVideo{})
	db.AutoMigrate(&models.FavouriteVideo{})
}
