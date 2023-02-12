package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 根据identity，查询关注信息
func FindUserFollowByIdentity(identity uint64) ([]models.UserFollow, error) {
	data := make([]models.UserFollow, 0)
	err := utils.DB.Where("user_identity = ? and status = 1", identity).Find(&data).Error
	return data, err
}

// 查询id对应的用户名
func FindUserName(userId string) (string, error) {
	user := models.UserBasic{}
	err := utils.DB.Where("identity = ?", userId).First(&user).Error
	return user.Username, err
}
