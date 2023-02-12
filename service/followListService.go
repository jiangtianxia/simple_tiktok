package service

import (
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 用户参数
type User struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

/**
 * @Author jiang
 * @Description 关注列表
 * @Date 12:00 2023/2/12
 **/
func FollowListService(c *gin.Context, userId uint64) ([]User, error) {
	// data := make([]User, 0)

	// 1、判断缓存中是否存在数据
	key := viper.GetString("redis.KeyFollowListPrefix") + strconv.Itoa(int(userId))

	if utils.RDB10.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库
		// 1）查询该用户所关注的用户数
		followerList, err := mysql.FindUserFollowByIdentity(userId)
		if err != nil {
			logger.SugarLogger.Error("FindUserFollowByIdentity Error：", err.Error())
			return nil, err
		}

		// 2）根据关注者的id，查询其信息
		for _, follower := range followerList {
			// 获取用户名
			_, err := GetUsername(c, strconv.Itoa(int(follower.UserIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetUsername Error：", err.Error())
				return nil, err
			}
		}
	}

	// 存在，则直接使用缓存
	return nil, nil
}

// 获取用户名
func GetUsername(c *gin.Context, identity string) (string, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyUserHashPrefix") + identity

	if utils.RDB1.Exists(c, key).Val() == 0 {
		// 存在，则获取，并返回
		userName, err := utils.RDB1.HGet(c, key, "username").Result()
		if err != nil {
			return "", err
		}
		return userName, nil
	}

	// 不存在，则在数据库中获取，同时将数据存入缓存
	username, err := mysql.FindUserName(identity)
	if err != nil {
		return "", err
	}

	return username, err
}
