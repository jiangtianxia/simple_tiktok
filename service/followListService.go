package service

import (
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 关注列表
 * @Date 12:00 2023/2/12
 **/
func FollowListService(c *gin.Context, userId uint64) ([]Author, error) {
	data := make([]Author, 0)

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
		pipeline := utils.RDB10.Pipeline()
		for _, follower := range followerList {
			// 获取用户名
			username, err := GetUsername(c, strconv.Itoa(int(follower.FollowerIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetUsername Error：", err.Error())
				return nil, err
			}

			// 获取关注数
			followCount, err := GetFollowCount(c, strconv.Itoa(int(follower.FollowerIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
				return nil, err
			}

			// 获取粉丝数
			followerCount, err := GetFollowerCount(c, strconv.Itoa(int(follower.FollowerIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
				return nil, err
			}

			// 判断是否互关
			// isFollow, err := IsFollow(c, strconv.Itoa(int(follower.UserIdentity)), strconv.Itoa(int(follower.FollowerIdentity)))
			// if err != nil {
			// 	logger.SugarLogger.Error("IsFollow Error：", err.Error())
			// 	return nil, err
			// }
			// 关注列表接口，因此已经关注该用户了
			userInfo := Author{
				Id:            follower.FollowerIdentity,
				Name:          username,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      true,
			}
			// 头部插入数据到key当中
			// LPUSH KEY_NAME VALUE1.. VALUEN
			pipeline.LPush(c, key, follower.FollowerIdentity)
			data = append(data, userInfo)
		}
		pipeline.Expire(c, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(c)
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			return nil, err
		}
		return data, nil
	}

	// 存在，则直接使用缓存
	// 获取列表中的全部元素
	// LRANGE KEY_NAME START END
	list, err := utils.RDB10.LRange(c, key, 0, -1).Result()
	if err != nil {
		logger.SugarLogger.Error("utils.RDB10.LRange Error：", err.Error())
		return nil, err
	}

	for _, followIdentity := range list {
		// 获取用户名
		username, err := GetUsername(c, followIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetUsername Error：", err.Error())
			return nil, err
		}

		// 获取关注数
		followCount, err := GetFollowCount(c, followIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
			return nil, err
		}

		// 获取粉丝数
		followerCount, err := GetFollowerCount(c, followIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
			return nil, err
		}

		// 判断是否互关
		// isFollow, err := IsFollow(c, strconv.Itoa(int(userId)), followIdentity)
		// if err != nil {
		// 	logger.SugarLogger.Error("IsFollow Error：", err.Error())
		// 	return nil, err
		// }

		id, _ := strconv.Atoi(followIdentity)
		userInfo := Author{
			Id:            uint64(id),
			Name:          username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		}
		data = append(data, userInfo)
	}
	return data, nil
}
