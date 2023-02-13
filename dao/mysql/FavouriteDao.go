package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

func IsFavourite(userId uint64, videoId string) bool {
	favourite := models.FavouriteVideo{}
	utils.DB.Where("user_identity=? and video_identity=?", userId, videoId).Find(&favourite)
	return favourite.Status == "1"
}

func UpdateFavourite(videoId string, userId uint64, change string) error {
	favourite := models.FavouriteVideo{}
	err := utils.DB.Model(&favourite).Where("user_identity=? and video_identity=?", userId, videoId).Updates(models.FavouriteVideo{Status: change}).Error
	return err
}

// 查询视频喜欢记录
func QueryVideoFavoriteCount(videoId *uint64) (*[]models.FavouriteVideo, error) {
	var videoFavoriteList []models.FavouriteVideo
	utils.DB.Table("favourite_video").Where("video_identity = ?", *videoId).Find(&videoFavoriteList)
	return &videoFavoriteList, nil
}

// 查询登录者是否喜欢该视频
func QueryIsFavorite(videoId *uint64, userId *uint64) (bool, error) {
	var videoFavorite models.FavouriteVideo
	result := utils.DB.Table("favourite_video").Where("video_identity = ? AND user_identity = ?", *videoId, *userId).Find(&videoFavorite)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
