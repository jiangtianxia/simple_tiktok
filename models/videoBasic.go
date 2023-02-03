package models

import (
	"gorm.io/gorm"
)

type VideoBasic struct {
	gorm.Model
	Identity     uint64 `gorm:"column:identity;type:int;"`           // 视频唯一标识
	UserIdentity uint64 `gorm:"column:user_identity;type:int;"`      // 作者ID
	PlayUrl      string `gorm:"column:play_url;type:varchar(255);"`  // 视频路径
	CoverUrl     string `gorm:"column:cover_url;type:varchar(255);"` // 封面路径
	Title        string `gorm:"column:title;type:text;"`             // 视频标题
	PublishTime  int64  `gorm:"column:publish_time;type:int;"`       // 发布时间
}

func (table *VideoBasic) TableName() string {
	return "video_basic"
}
