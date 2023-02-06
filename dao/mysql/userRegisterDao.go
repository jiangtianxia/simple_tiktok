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
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

//查询是否存在
func IsExist(username string) bool {
	var user models.UserBasic
	utils.DB.Select("username").Where("username = ?", username).First(&user)
	return user.Identity != 0
}

// 用户增加
func AddUserBasic(userbasic *models.UserBasic) error {
	if userbasic == nil {
		return errors.New("空指针")
	}
	utils.DB.Create(&userbasic)
	return nil
}