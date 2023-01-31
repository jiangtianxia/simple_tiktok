package mysql

import (
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 根据identity查询用户信息
func FindUserByIdentity(identity uint64) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if err := utils.DB.Table("user_basic").Where("identity = ?", identity).First(&user).Error; err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	} else {
		return &user, nil
	}
}
