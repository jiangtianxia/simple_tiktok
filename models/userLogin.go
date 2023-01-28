package models

import "gorm.io/gorm"

type UserLogin struct {
	gorm.Model
	Identity uint64 `gorm:"column:identity;type:int;"`         // 用户唯一标识
	Username string `gorm:"column:username;type:varchar(36);"` // 用户名
	Password string `gorm:"column:password;type:varchar(36);"` // 密码
}

var userLogin UserLogin

func IsExist(username string) error {
	return DB.Where("username = ?", username).First(&userLogin).Error
}

func QueryPassword(username string) string {
	return DB.Select("password").Where("username=?", username).First(&userLogin)
}

func QueryIdentity(username string) string {
	return DB.Select("identity").Where("username=?", username).First(&userLogin)
}
