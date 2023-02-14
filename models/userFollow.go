package models

import (
	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	UserIdentity     uint64 `gorm:"column:user_identity;type:int;"`     // 用户ID
	FollowerIdentity uint64 `gorm:"column:follower_identity;type:int;"` // 关注者ID
	Status           int    `gorm:"column:status;type:tinyint(1);"`     // 关注状态（0表示未关注，1表示已关注）
}

func (table *UserFollow) TableName() string {
	return "user_follow"
}
