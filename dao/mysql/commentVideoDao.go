package mysql

import (
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
