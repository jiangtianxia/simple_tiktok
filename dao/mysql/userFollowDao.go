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

// 根据identity查询粉丝信息
func FindFollower(identity string) ([]models.UserFollow, error) {
	data := make([]models.UserFollow, 0)
	err := utils.DB.Where("follower_identity = ? and status = 1", identity).Find(&data).Error
	return data, err
}

// 判断用户是否互相关注
func IsFollow(identity string, follower string) (int64, error) {
	var cnt int64
	err := utils.DB.Model(new(models.UserFollow)).
		Where("user_identity = ? and follower_identity = ? and status = 1", identity, follower).
		Count(&cnt).Error
	return cnt, err
}
