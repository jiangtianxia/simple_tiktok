package models

import (
	"gorm.io/gorm"
)

type FavouriteVideo struct {
	gorm.Model
	Identity      uint64 `gorm:"column:identity;type:int;"`       // 点赞关系表唯一标识
	VideoIdentity uint64 `gorm:"column:video_identity;type:int;"` // 点赞视频id
	UserIdentity  uint64 `gorm:"column:user_identity;type:int;"`  // 点赞用户id
	Status        string `gorm:"column:text;type:tinyint(1);"`    // 评论状态（是否点赞,0表示未点赞，1表示点赞）
}

func (table *FavouriteVideo) TableName() string {
	return "favourite_video"
}
