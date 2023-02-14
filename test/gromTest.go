package main

import (
	"simple_tiktok/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main1() {
	db, err := gorm.Open(mysql.Open("test:674092@tcp(101.43.157.116:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.VideoBasic{})
	db.AutoMigrate(&models.CommentVideo{})
	db.AutoMigrate(&models.FavouriteVideo{})
	db.AutoMigrate(&models.UserFollow{})
	db.AutoMigrate(&models.UserMessage{})
}
