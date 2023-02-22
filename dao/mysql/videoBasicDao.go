package mysql

import (
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"time"

	"github.com/spf13/viper"
)

// 将数据插入videoBasic
func CreateVideoBasic(videoInfo models.VideoBasic) error {
	return utils.DB.Create(&videoInfo).Error
}

// 根据identity查询用户上传的视频列表
func FindVideoByUserIdentity(userid uint64) ([]models.VideoBasic, error) {
	videoList := make([]models.VideoBasic, 0)
	err := utils.DB.Table("video_basic").Where("user_identity = ?", userid).Order("publish_time asc").Find(&videoList).Error
	return videoList, err
}

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

// 查询视频列表
func QueryVideoList(userId *uint64) (*[]models.VideoBasic, error) {
	var videoInfoList []models.VideoBasic
	utils.DB.Table("video_basic").Where("user_identity = ?", *userId).Find(&videoInfoList)
	return &videoInfoList, nil
}

// 查询视频信息
func QueryVideoInfoByVideoId(videoId *uint64) (*models.VideoBasic, error) {
	var videoInfo models.VideoBasic
	utils.DB.Table("video_basic").Where("identity = ?", *videoId).Find(&videoInfo)
	return &videoInfo, nil
}

// 获取点赞信息
func QueryFavoriteHistoryByUserId(userId *uint64) (*[]models.FavouriteVideo, error) {
	var favoriteList []models.FavouriteVideo
	utils.DB.Table("favourite_video").Where("user_identity = ?", *userId).Find(&favoriteList)
	return &favoriteList, nil
}

// 获取作品数
func GetWorkCount(identity uint64) ([]models.VideoBasic, error) {
	var VideoList []models.VideoBasic
	err := utils.DB.Model(new(models.VideoBasic)).Where("user_identity = ?", identity).Find(&VideoList).Error
	return VideoList, err
}
