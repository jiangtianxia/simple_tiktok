package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

/** 
 * 获取发布列表的数据库查询操作
 * @Author pjh
 * @Summary 
 * @Tags 
 **/

// 查询视频列表
func QueryVideoList(userId *uint64) (*[]models.VideoBasic, error) {
	var videoInfoList []models.VideoBasic
	utils.DB.Table("video_basic").Where("user_identity = ?", *userId).Find(&videoInfoList)
	return &videoInfoList, nil
}

// 查询id对应的用户名
func QueryAuthorName(userId *uint64) (*string, error) {
	var author models.UserBasic
	utils.DB.Table("user_basic").Where("identity = ?", *userId).Find(&author)
	return &author.Username, nil
}

// 查询视频喜欢记录
func QueryVideoFavoriteCount (videoId *uint64) (*[]models.FavouriteVideo, error) {
	var videoFavoriteList []models.FavouriteVideo
	utils.DB.Table("video_favorite").Where("video_identity = ?", *videoId).Find(&videoFavoriteList)
	return &videoFavoriteList, nil
}

// 查询视频的评论数
func QueryCommentCount (videoId *uint64) (*int64, error) {
	var commentVideo models.CommentVideo
	result := utils.DB.Table("comment_video").Where("video_identity = ?", *videoId).Find(&commentVideo)
	return &result.RowsAffected, nil
}

// 查询登录者是否喜欢该视频
func QueryIsFavorite(videoId *uint64, userId *uint64) (bool, error) {
	var videoFavorite models.FavouriteVideo
	result := utils.DB.Table("video_favorite").Where("video_identity = ? AND user_identity = ?", *videoId, *userId).Find(&videoFavorite)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// 查询视频信息
func QueryVideoInfoByVideoId(videoId *uint64) (*models.VideoBasic, error) {
	var videoInfo models.VideoBasic
	utils.DB.Table("video_basic").Where("identity = ?", *videoId).Find(&videoInfo)
	return &videoInfo, nil
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
