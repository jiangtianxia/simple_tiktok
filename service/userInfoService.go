package service

import (
	"encoding/json"
	"fmt"
	"simple_tiktok/dao/mysql"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author Xiaoyu Zhang
 * @Description 获取用户信息
 * @Date 14:00 2023/1/31
 **/
func UserInfo(c *gin.Context, userId string) (map[string]interface{}, error) {
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + userId
	// 判断是否有缓存
	if utils.RDB1.Exists(c, hashKey).Val() == 0 {
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
			redisErr := myRedis.RedisAddUserInfo(hashKey, map[string]interface{}{
				"identity": viper.GetInt("redis.defaultErrorIdentity"),
			})
			if redisErr != nil {
				logger.SugarLogger.Error(err)
				return nil, redisErr
			}
			return nil, err
		}

		// 使用时间队列新增缓存
		var userInfo = models.UserBasic{
			Identity: user.Identity,
			Username: user.Username,
		}
		data, _ := json.Marshal(userInfo)
		redisTopic := viper.GetString("rocketmq.redisTopic")
		Producer := viper.GetString("rocketmq.redisProducer")
		tag := viper.GetString("rocketmq.userInfoTag")
		msg, err := utils.SendMsg(c, Producer, redisTopic, tag, data)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg)
		fmt.Println("新增缓存：", userInfo)
		//// 新增缓存
		//var newCathe = map[string]interface{}{
		//	"identity": user.Identity,
		//	"username": user.Username,
		//}
		//
		//err = myRedis.RedisAddUserInfo(hashKey, newCathe)
		//if err != nil {
		//	logger.SugarLogger.Error(err)
		//	return nil, err
		//}

		// 返回结果
		var res = map[string]interface{}{
			"identity":       user.Identity,
			"username":       user.Username,
			"follow_count":   0,
			"follower_count": 0,
			"is_follow":      false,
		}
		return res, nil
	}
	// 使用缓存
	cathe := utils.RDB1.HGetAll(c, hashKey).Val()
	identity, _ := strconv.Atoi(cathe["identity"])
	var res = map[string]interface{}{
		"identity":       identity,
		"username":       cathe["username"],
		"follow_count":   0,
		"follower_count": 0,
		"is_follow":      false,
	}
	fmt.Println("已有缓存", cathe)

	return res, nil
}
