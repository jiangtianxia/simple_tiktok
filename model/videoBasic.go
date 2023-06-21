package model

import (
	"gorm.io/gorm"
)

type VideoBasic struct {
	gorm.Model
	AuthorId    uint   `gorm:"column:author_id;type:uint;comment:'作者ID'"`         // 作者ID
	PlayUrl     string `gorm:"column:play_url;type:varchar(255);comment:'视频路径'"`  // 视频路径
	CoverUrl    string `gorm:"column:cover_url;type:varchar(255);comment:'封面路径'"` // 封面路径
	Title       string `gorm:"column:title;type:text;comment:'视频标题'"`             // 视频标题
	PublishTime int64  `gorm:"column:publish_time;type:int(64);comment:'发布时间'"`   // 发布时间
}

func (table *VideoBasic) TableName() string {
	return "video_basic"
}
