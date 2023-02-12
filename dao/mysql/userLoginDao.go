package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

func IsExist(username string) error {
	userLogin := models.UserBasic{}
	return utils.DB.Where("username = ?", username).First(&userLogin).Error
}

func QueryInfo(username string) models.UserBasic {
	userLogin := models.UserBasic{}
	utils.DB.Where("username=?", username).First(&userLogin)
	return userLogin
}
