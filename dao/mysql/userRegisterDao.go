package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

//请求体
type RegisterRequire struct {
	Username string 
	Password    string 
}

type RegisterResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}

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