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
	"simple_tiktok/utils"
	"sync"
)

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
	utils.DB.Where("username = ?", username).First(&user)
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
	utils.DB.Select(&models.UserBasic{}).Where("identity = ?", userId)
	if userbasic.Identity == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}

func (u *UserBasicDAO) AddUserBasic(userbasic *models.UserBasic) error {
	if userbasic == nil {
		return errors.New("空指针")
	}

	return utils.DB.Create(userbasic).Error
}

func (u *UserBasicDAO) IsExist(id int64) bool {
	var userbasic models.UserBasic
	if err := utils.DB.Where("identity = ?", id).Select("identity").First(&userbasic).Error; err != nil {
		log.Println(err)
	}
	if userbasic.Identity == 0 {
		return false
	}
	return true
}
