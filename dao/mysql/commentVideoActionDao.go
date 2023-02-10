package mysql

import (
	"errors"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// Comment
// 评论信息-数据库中的结构体-dao层使用
// type CommentVideo struct {
// 	gorm.Model
// 	Identity      uint64 `gorm:"column:identity;type:int;"`            // 评论唯一标识
// 	VideoIdentity uint64 `gorm:"column:video_identity;type:int;"`      // 视频ID
// 	UserIdentity  uint64 `gorm:"column:user_identity;type:int;"`       // 用户ID
// 	Text          string `gorm:"column:text;type:text;"`               // 评论内容
// 	CommentTime   string `gorm:"column:comment_time;type:varchar(10)"` // 评论时间
// }

// func (table *CommentVideo) TableName() string {
// 	return "comment_video"
// }

// 发表评论
func AddComment(comment models.CommentVideo) error {
	err := utils.DB.Model(models.CommentVideo{}).Create(&comment).Error
	if err != nil {
		return errors.New("发表评论失败")
	}
	return nil
}

// 删除评论，传入评论id
func DelComment(identity uint64) error {
	var commentInfo models.CommentVideo
	result := utils.DB.Model(models.CommentVideo{}).Where("identity = ?", identity).First(&commentInfo)
	if result.RowsAffected == 0 {
		return errors.New("该评论不存在")
	}

	err := utils.DB.Delete(commentInfo).Error
	if err != nil {
		return errors.New("删除评论失败")
	}
	return nil
}
