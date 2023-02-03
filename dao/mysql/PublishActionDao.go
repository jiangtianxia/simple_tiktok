package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
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
