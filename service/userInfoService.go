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
			redisErr := myRedis.RedisAddUserInfo(c, hashKey, map[string]interface{}{
				"identity": viper.GetInt("redis.defaultErrorIdentity"),
			})
			if redisErr != nil {
				logger.SugarLogger.Error(err)
				return nil, redisErr
			}
			return nil, err
		}

		// 新增缓存
		var newCathe = map[string]interface{}{
			"identity": user.Identity,
			"username": user.Username,
		}

		err = myRedis.RedisAddUserInfo(c, hashKey, newCathe)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		fmt.Println("新增缓存：", newCathe)

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

	//// 使用时间队列
	//var userInfo = models.UserBasic{
	//	Identity: 1,
	//	Username: "merry",
	//	Password: "123456",
	//}
	//fmt.Println(userInfo)
	//redisTopic := viper.GetString("rocketmq.redisTopic")
	//Producer := viper.GetString("rocketmq.redisProducer")
	//tag := viper.GetString("rocketmq.userInfoTag")
	//data, _ := json.Marshal(userInfo)
	//msg, err := rocket.SendMsg(c, Producer, redisTopic, tag, data)
	//fmt.Println(msg)
	//if err != nil {
	//	return nil, err
	//}

	return res, nil
}
