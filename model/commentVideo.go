package model

import (
	"gorm.io/gorm"
)

type CommentVideo struct {
	gorm.Model
	VideoId     uint   `gorm:"column:video_id;type:uint;comment:'视频ID'"`                   // 视频ID
	UserId      uint   `gorm:"column:user_id;type:uint;comment:'用户ID'"`                    // 用户ID
	Text        string `gorm:"column:text;type:text;comment:'评论内容'"`                       // 评论内容
	CommentTime string `gorm:"column:comment_time;type:varchar(10);comment:'评论时间, MM-DD'"` // 评论时间
}

func (table *CommentVideo) TableName() string {
	return "comment_video"
}
