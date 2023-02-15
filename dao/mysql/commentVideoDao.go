package mysql

import (
	"errors"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 查询视频的评论数
func QueryCommentCount(videoId *uint64) (*int64, error) {
	var commentVideo models.CommentVideo
	result := utils.DB.Table("comment_video").Where("video_identity = ?", *videoId).Find(&commentVideo)
	return &result.RowsAffected, nil
}

// 查询评论信息
func QueryVideoCommentInfo(videoId *uint64) (*[]models.CommentVideo, error) {
	var commentList []models.CommentVideo
	utils.DB.Table("comment_video").Where("video_identity = ?", *videoId).Find(&commentList)
	return &commentList, nil
}

// 获取点赞信息
func QueryFavoriteInfo(videoId *uint64) (*[]models.FavouriteVideo, error) {
	var favoriteList []models.FavouriteVideo
	utils.DB.Table("comment_video").Where("video_identity = ?", *videoId).Find(&favoriteList)
	return &favoriteList, nil
}

// 发表评论，传入评论结构体
func AddComment(comment models.CommentVideo) error {
	err := utils.DB.Model(models.CommentVideo{}).Create(&comment).Error
	if err != nil {
		return errors.New("发表评论失败")
	}
	return nil
}

// 删除评论，传入评论id
func DelComment(identity uint64) error {
	err := utils.DB.Model(models.CommentVideo{}).Delete("identity = ?", identity).Error

	if err != nil {
		return errors.New("删除评论失败")
	}
	return nil
}

// 根据视频id查询视频作者id
func SearchAuthorIdByVideoId(identity uint64) (uint64, error) {
	var commentInfo models.VideoBasic
	result := utils.DB.Model(models.VideoBasic{}).Where("identity = ?", identity).First(&commentInfo)
	if result.RowsAffected == 0 {
		return 0, errors.New("该视频作者不存在")
	}
	return commentInfo.UserIdentity, nil
}
