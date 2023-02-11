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
