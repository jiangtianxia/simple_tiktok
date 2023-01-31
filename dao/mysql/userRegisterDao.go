/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 10:22:09
 * @LastEditTime: 2023-01-29 13:46:00
 * @LastEditors:
 */
package mysql

import (
	"errors"
	"log"
	"simple_tiktok/models"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dsn := "user:test@tcp(127.0.0.1:8080)/init?charset=utf8mb4&parseTime=True&loc=Local"
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.UserBasic{})
}

type UserRegisterDAO struct {
}

var (
	userRegisterDao  *UserRegisterDAO
	userRegisterOnce sync.Once
)

func NewUserRegisterDAO() *UserRegisterDAO {
	userRegisterOnce.Do(func() {
		userRegisterDao = new(UserRegisterDAO)
	})
	return userRegisterDao
}

func (req *UserRegisterDAO) IsExist(username string) bool {
	var user models.UserBasic
	db.Where("username = ?", username).First(&user)
	return user.Identity != 0
}

// 用户增加
type UserBasicDAO struct {
}

var (
	userBasicDAO  *UserBasicDAO
	userBasicOnce sync.Once
)

func NewUserBasicDAO() *UserBasicDAO {
	userBasicOnce.Do(func() {
		userBasicDAO = new(UserBasicDAO)
	})
	return userBasicDAO
}

func (u *UserBasicDAO) QueryUserBasicById(userId int64, userbasic *models.UserBasic) error {
	if userbasic == nil {
		return errors.New("空指针")
	}
	db.Select(&models.UserBasic{}).Where("identity = ?", userId)
	if userbasic.Identity == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}

func (u *UserBasicDAO) AddUserBasic(userbasic *models.UserBasic) error {
	if userbasic == nil {
		return errors.New("空指针")
	}

	return db.Create(userbasic).Error
}

func (u *UserBasicDAO) IsExist(id int64) bool {
	var userbasic models.UserBasic
	if err := db.Where("identity = ?", id).Select("identity").First(&userbasic).Error; err != nil {
		log.Println(err)
	}
	if userbasic.Identity == 0 {
		return false
	}
	return true
}
