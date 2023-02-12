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
