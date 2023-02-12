/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 10:22:09
 * @LastEditTime: 2023-01-29 13:46:00
 * @LastEditors:
 */
package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

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