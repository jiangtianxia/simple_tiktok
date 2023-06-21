package model

import (
	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	UserId     uint `gorm:"column:user_id;type:uint;comment:'用户ID'"`                      // 用户ID
	FollowerId uint `gorm:"column:follower_id;type:uint;comment:'关注者ID'"`                 // 关注者ID
	Status     int  `gorm:"column:status;type:tinyint(1);comment:'关注状态(0表示未关注, 1表示已关注)'"` // 关注状态(0表示未关注, 1表示已关注)
}

func (table *UserFollow) TableName() string {
	return "user_follow"
}
