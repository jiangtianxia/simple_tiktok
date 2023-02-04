package models

import (
	"gorm.io/gorm"
)

type userFollow struct {
	gorm.Model
	UserIdentity   uint64 `gorm:"column:user_identity;type:int;"`   // 用户ID
	FollowIdentity uint64 `gorm:"column:follow_identity;type:int;"` // 关注者ID
	Status         string `gorm:"column:status;type:tinyint(1);"`   // 关注状态（0表示未关注，1表示已关注）
}

func (table *userFollow) TableName() string {
	return "user_follow"
}
