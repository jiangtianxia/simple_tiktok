package mysql

import (
	"github.com/spf13/viper"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"time"
)

func FindVideosByLatestTime(latestTime time.Time) ([]models.VideoBasic, error) {
	videos := make([]models.VideoBasic, viper.GetInt("feedVideoCnt"))
	result := utils.DB.Table("video_basic").Where("publish_time<=?", latestTime).Order("publish_time desc").Limit(viper.GetInt("feedVideoCnt")).Find(&videos)
	if result.Error != nil {
		logger.SugarLogger.Error(result.Error)
		return videos, result.Error
	}
	return videos, nil
}

// 根据identity查询视频信息
func FindVideoById(identity uint64) (*models.VideoBasic, error) {
	video := models.VideoBasic{}
	if err := utils.DB.Table("video_basic").Where("identity = ?", identity).First(&video).Error; err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	} else {
		return &video, nil
	}
}

// 查询视频Identity集合
func FindVideoIdentityByLatestTime(latestTime time.Time) (*[]uint64, error) {
	var Ids []uint64
	if err := utils.DB.Table("video_basic").Where("publish_time <= ?", latestTime).Pluck("identity", &Ids).Limit(viper.GetInt("feedVideoCnt")).Error; err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	} else {
		return &Ids, nil
	}
}
