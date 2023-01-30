package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

func IsExist(username string) error {
	userLogin := models.UserBasic{}
	return utils.DB.Where("username = ?", username).First(&userLogin).Error
}

func QueryPassword(username string) string {
	userLogin := models.UserBasic{}
	utils.DB.Select("password").Where("username=?", username).First(&userLogin)
	return userLogin.Password
}

func QueryIdentity(username string) uint64 {
	userLogin := models.UserBasic{}
	utils.DB.Select("identity").Where("username=?", username).First(&userLogin)
	return userLogin.Identity
}
