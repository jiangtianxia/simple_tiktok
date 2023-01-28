package models

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity uint64 `gorm:"column:identity;type:uint64;"`      // 用户唯一标识
	Username string `gorm:"column:username;type:varchar(36);"` // 用户名
	Password string `gorm:"column:password;type:varchar(36);"` // 密码
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
