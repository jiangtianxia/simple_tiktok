package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 将数据插入videoBasic
func CreateVideoBasic(videoInfo models.VideoBasic) error {
	return utils.DB.Create(&videoInfo).Error
}
