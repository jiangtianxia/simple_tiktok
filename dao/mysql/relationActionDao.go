package mysql

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 根据identity查询用户信息
func FindUserByIdentityCount(identity uint64) (int64, error) {
	var cnt int64
	err := utils.DB.Model(new(models.UserBasic)).Where("identity = ?", identity).Count(&cnt).Error
	return cnt, err
}

// 查询关注表中是否存在数据
func FindUserFollowByIdentityCount(UserIdentity uint64, followIdentity uint64) (int64, error) {
	var cnt int64
	err := utils.DB.Model(new(models.UserFollow)).Where("user_identity = ? and follower_identity = ?", UserIdentity, followIdentity).Count(&cnt).Error
	return cnt, err
}

// 在关注表中插入数据
func CreateUserFollow(info models.UserFollow) error {
	return utils.DB.Create(&info).Error
}

// 修改关注表中的数据
func UpdateUserFollow(info models.UserFollow) error {
	return utils.DB.Model(new(models.UserFollow)).
		Where("user_identity = ? and follower_identity = ?", info.UserIdentity, info.FollowerIdentity).
		Update("status", info.Status).Error
}
