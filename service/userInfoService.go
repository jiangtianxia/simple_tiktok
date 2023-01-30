package service

import (
	"fmt"
	"simple_tiktok/dao/mysql"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func UserInfo(c *gin.Context, userId string) (interface{}, error) {
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + userId
	// 判断是否有缓存
	if utils.RDB0.Exists(c, hashKey).Val() == 0 {
		// 查询数据库
		idNum, err := strconv.Atoi(userId)
		identityUint64 := uint64(idNum)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		user, err := mysql.FindUserByIdentity(identityUint64)
		if err != nil {
			// 防止缓存击穿
			redisErr := myRedis.RedisAddUserInfo(c, hashKey, map[string]interface{}{
				"identity": -1,
			})
			if redisErr != nil {
				logger.SugarLogger.Error(err)
				return nil, redisErr
			}
			return nil, err
		}
		var res = map[string]interface{}{
			"identity":       user.Identity,
			"username":       user.Username,
			"follow_count":   0,
			"follower_count": 0,
			"is_follow":      false,
		}

		// 新增缓存
		err = myRedis.RedisAddUserInfo(c, hashKey, res)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}

		fmt.Println("数据库")
		return res, nil

	}
	// 使用缓存
	cathe := utils.RDB0.HGetAll(c, hashKey).Val()
	identity, _ := strconv.Atoi(cathe["identity"])
	followCount, _ := strconv.Atoi(cathe["follow_count"])
	followerCount, _ := strconv.Atoi(cathe["follower_count"])
	isFollow, _ := strconv.ParseBool(cathe["is_follow"])
	var res = map[string]interface{}{
		"identity":       identity,
		"username":       cathe["username"],
		"follow_count":   followCount,
		"follower_count": followerCount,
		"is_follow":      isFollow,
	}
	fmt.Println("缓存")
	return res, nil
}
