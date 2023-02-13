package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: dao层-使用数据库中的结构体，查询用户是否存在，添加用户信息
 **/

// 查询是否存在
func IsExist(username string) bool {
	var user models.UserBasic
	var size int64
	utils.DB.Model(&models.UserBasic{}).Select("username").Where("username = ?", username).Scan(&user).Count(&size)
	return size != 0
}

// 用户增加
func AddUserBasic(userbasic *models.UserBasic) error {
	return utils.DB.Create(&userbasic).Error
}