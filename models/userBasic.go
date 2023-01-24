package models

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity uint64 `gorm:"coulmn:identity;type:int"`         // 用户唯一标识
	Username string `gorm:"coulmn:username;type:varchar(36)"` // 用户名
	Password string `gorm:"coulmn:password;type:varchar(36)"` // 密码
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
