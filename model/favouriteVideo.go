package model

import (
	"gorm.io/gorm"
)

type FavouriteVideo struct {
	gorm.Model
	VideoId uint `gorm:"column:video_id;type:uint;comment:'点赞视频id'"`                   // 点赞视频id
	UserId  uint `gorm:"column:user_id;type:uint;comment:'点赞用户id'"`                    // 点赞用户id
	Status  int  `gorm:"column:status;type:tinyint(1);comment:'点赞状态(0表示未点赞, 1表示已点赞)'"` // 点赞状态(0表示未点赞, 1表示已点赞)
}

func (table *FavouriteVideo) TableName() string {
	return "favourite_video"
}
