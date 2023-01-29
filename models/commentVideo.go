package models

import (
	"time"

	"gorm.io/gorm"
)

type CommentVideo struct {
	gorm.Model
	Identity      uint64    `gorm:"column:identity;type:int(64);"`       // 评论唯一标识
	VideoIdentity uint64    `gorm:"column:video_identity;type:int(64);"` // 视频ID
	UserIdentity  uint64    `gorm:"column:user_identity;type:int(64);"`  // 用户ID
	Text          string    `gorm:"column:text;type:varchar(100);"`      // 评论内容
	CommentTime   time.Time `gorm:"column:comment_time;type:datetime"`   // 评论时间
}

func (table *CommentVideo) TableName() string {
	return "comment_video"
}
