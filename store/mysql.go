package store

import (
	"log"
	"os"
	"simple_tiktok/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

// 初始化数据库
func InitDB(conn string) (err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值，超过1秒的sql查询会被记录到日志中
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	db, err = gorm.Open(mysql.Open(conn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// 设置最大连接数
	sqlDB.SetMaxIdleConns(1000)
	// 设置最大的空闲连接数
	sqlDB.SetMaxOpenConns(100000)
	return
}

// 自动迁移表
func AutoMigrate() error {
	if err := db.AutoMigrate(&model.UserBasic{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.VideoBasic{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.FavouriteVideo{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.CommentVideo{}); err != nil {
		return err
	}

	return nil
}
