/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 10:22:09
 * @LastEditTime: 2023-01-29 13:46:00
 * @LastEditors:
 */
package mysql

import (
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 查询是否存在
func UserIsExist(username string) bool {
	var user models.UserBasic
	var size int64
	utils.DB.Model(&models.UserBasic{}).Select("username").Where("username = ?", username).Scan(&user).Count(&size)
	return size != 0
}

// 用户增加
func AddUserBasic(userbasic models.UserBasic) error {
	return utils.DB.Create(&userbasic).Error
}

// 查询用户信息
func QueryUserInfo(username string) models.UserBasic {
	userLogin := models.UserBasic{}
	utils.DB.Where("username=?", username).First(&userLogin)
	return userLogin
}

// 根据identity查询用户信息
func FindUserByIdentity(identity uint64) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if err := utils.DB.Table("user_basic").Where("identity = ?", identity).First(&user).Error; err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	}
	return &user, nil
}

// 查询id对应的用户名
func FindUserName(userId string) (string, error) {
	user := models.UserBasic{}
	err := utils.DB.Where("identity = ?", userId).First(&user).Error
	return user.Username, err
}

// 查询id对应的用户名
func QueryAuthorName(userId *uint64) (*string, error) {
	var author models.UserBasic
	utils.DB.Table("user_basic").Where("identity = ?", *userId).Find(&author)
	return &author.Username, nil
}

// 根据identity查询用户信息
func FindUserByIdentityCount(identity uint64) (int64, error) {
	var cnt int64
	err := utils.DB.Model(new(models.UserBasic)).Where("identity = ?", identity).Count(&cnt).Error
	return cnt, err
}
