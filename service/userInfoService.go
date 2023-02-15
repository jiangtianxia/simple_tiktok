package service

import (
	"errors"
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
func UserInfo(c *gin.Context, loginUser uint64, userId string) (Author, error) {
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + userId
	// 判断是否有缓存
	if utils.RDB1.Exists(c, hashKey).Val() == 0 {
		// 查询数据库
		idNum, err := strconv.Atoi(userId)
		identityUint64 := uint64(idNum)
		if err != nil {
			logger.SugarLogger.Error(err)
			return Author{}, err
		}
		user, err := mysql.FindUserByIdentity(identityUint64)
		if err != nil {
			// 防止缓存击穿
			redisErr := myRedis.RedisAddUserInfo(c, hashKey, map[string]interface{}{
				"identity": viper.GetInt("redis.defaultErrorIdentity"),
			})
			if redisErr != nil {
				logger.SugarLogger.Error(err)
				return Author{}, redisErr
			}
			return Author{}, err
		}

		// 新增缓存
		var newCathe = map[string]interface{}{
			"identity": user.Identity,
			"username": user.Username,
		}

		err = myRedis.RedisAddUserInfo(c, hashKey, newCathe)
		if err != nil {
			logger.SugarLogger.Error(err)
		}

		// 获取关注数
		followCount, err := GetFollowCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
			return Author{}, err
		}

		// 获取粉丝数
		followerCount, err := GetFollowerCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
			return Author{}, err
		}

		// 判断是否关注用户
		flag := false
		if user.Identity == loginUser {
			flag = true
		} else {
			flag, err = IsFollow(c, strconv.Itoa(int(user.Identity)), strconv.Itoa(int(loginUser)))
			if err != nil {
				logger.SugarLogger.Error("IsFollow Error：", err.Error())
				return Author{}, err
			}
		}

		// 返回结果
		res := Author{
			Id:            user.Identity,
			Name:          user.Username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      flag,
		}
		return res, nil
	}
	// 使用缓存
	cathe := utils.RDB1.HGetAll(c, hashKey).Val()
	identity, _ := strconv.Atoi(cathe["identity"])
	if identity == viper.GetInt("redis.defaultErrorIdentity") {
		return Author{}, errors.New("用户不存在")
	}

	// 获取关注数
	followCount, err := GetFollowCount(c, userId)
	if err != nil {
		logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
		return Author{}, err
	}

	// 获取粉丝数
	followerCount, err := GetFollowerCount(c, userId)
	if err != nil {
		logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
		return Author{}, err
	}

	// 判断是否关注该用户
	flag := false
	if uint64(identity) == loginUser {
		flag = true
	} else {
		flag, err = IsFollow(c, userId, strconv.Itoa(int(loginUser)))
		// fmt.Println(flag)
		if err != nil {
			logger.SugarLogger.Error("IsFollow Error：", err.Error())
			return Author{}, err
		}
	}

	res := Author{
		Id:            uint64(identity),
		Name:          cathe["username"],
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      flag,
	}
	return res, nil
}
