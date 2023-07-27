package model

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(36);index:username_index;comment:'用户名, 最长32个字符'"` // 用户名, 最长32个字符
	Password string `gorm:"column:password;type:varchar(36);comment:'密码, 最长32个字符'"`                       // 密码, 最长32个字符
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
