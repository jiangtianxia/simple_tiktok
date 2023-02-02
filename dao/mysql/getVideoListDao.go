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
	utils.DB.Table("video_basic").Where("user_identity = ?", userId).Find(&videoInfoList)
	return &videoInfoList, nil
}
// 查询id对应的用户名
func QueryAuthorName(userId *uint64) (*string, error) {
	var author models.UserBasic
	utils.DB.Table("user_basic").Where("user_identity = ?", userId).Find(&author)
	return &author.Username, nil
}
// 查询视频喜欢数
func QueryVideoFavoriteCount (videoId *uint64) (*int64, error) {
	var videoFavorite models.FavouriteVideo
	result := utils.DB.Table("video_favorite").Where("video_identity = ?", videoId).Find(&videoFavorite)
	return &result.RowsAffected, nil
}
// 查询视频的评论数
func QueryCommentCount (videoId *uint64) (*int64, error) {
	var commentVideo models.CommentVideo
	result := utils.DB.Table("comment_video").Where("video_identity = ?", videoId).Find(&commentVideo)
	return &result.RowsAffected, nil
}
// 查询登录者是否喜欢该视频
func QueryIsFavorite(videoId *uint64, userId *uint64) (bool, error) {
	var videoFavorite models.FavouriteVideo
	result := utils.DB.Table("video_favorite").Where("video_identity = ? AND user_identity = ?", videoId, userId).Find(&videoFavorite)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
