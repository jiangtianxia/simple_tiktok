package models

import (
	"gorm.io/gorm"
)

type userMessage struct {
	gorm.Model
	Identity         uint64 `gorm:"column:identity;type:int;"`           // 聊天记录一标识
	ToUserIdentity   uint64 `gorm:"column:to_user_identity;type:int;"`   // 接收者ID
	FromUserIdentity uint64 `gorm:"column:from_user_identity;type:int;"` // 发送者ID
	Content          string `gorm:"column:text;type:text;"`              // 消息内容
	CreateTime       int64  `gorm:"column:create_time;type:int"`         //发送时间
}

func (table *userMessage) TableName() string {
	return "user_message"
}
