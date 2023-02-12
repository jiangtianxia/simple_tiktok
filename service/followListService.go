package service

import (
	"simple_tiktok/dao/mysql"
	Myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
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
	data := make([]User, 0)

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
			userInfo := User{
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
		userInfo := User{
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

// 获取用户名
func GetUsername(c *gin.Context, identity string) (string, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyUserHashPrefix") + identity

	if utils.RDB1.Exists(c, key).Val() == 0 {
		// 不存在，则在数据库中获取，同时将数据存入缓存
		username, err := mysql.FindUserName(identity)
		if err != nil {
			return "", err
		}

		err = Myredis.RedisAddUserInfoHash(c, key, map[string]interface{}{
			"identity": identity,
			"username": username,
		})
		if err != nil {
			return "", err
		}
		return username, err
	}

	// 存在，则获取，并返回
	userName, err := utils.RDB1.HGet(c, key, "username").Result()
	if err != nil {
		return "", err
	}
	return userName, nil
}

// 获取关注数
func GetFollowCount(c *gin.Context, identity string) (int64, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowListPrefix") + identity

	if utils.RDB10.Exists(c, key).Val() == 0 {
		// 不存在，则在数据库中获取，同时将数据存入缓存
		user_id, _ := strconv.Atoi(identity)
		followerList, err := mysql.FindUserFollowByIdentity(uint64(user_id))
		if err != nil {
			logger.SugarLogger.Error("FindUserFollowByIdentity Error：", err.Error())
			return 0, err
		}

		// 将数据存入缓存
		// 头部插入数据到key当中
		// LPUSH KEY_NAME VALUE1.. VALUEN
		pipeline := utils.RDB10.Pipeline()
		for _, follower := range followerList {
			pipeline.LPush(c, key, follower.FollowerIdentity)
		}
		pipeline.Expire(c, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(c)
		return int64(len(followerList)), err
	}

	// 存在，则获取，并返回
	return utils.RDB10.LLen(c, key).Result()
}

// 获取粉丝数
func GetFollowerCount(c *gin.Context, identity string) (int64, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + identity

	if utils.RDB11.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库，将数据存入缓存，并返回
		followerList, err := mysql.FindFollower(identity)
		if err != nil {
			return 0, nil
		}

		//  ZADD KEY_NAME SCORE1 VALUE1.. SCOREN VALUEN
		pipeline := utils.RDB11.Pipeline()
		for _, follower := range followerList {
			pipeline.ZAdd(c, key, redis.Z{
				Member: follower.UserIdentity,
				Score:  1,
			})
		}
		pipeline.Expire(c, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(c)
		return int64(len(followerList)), err
	}

	// 存在，则获取，并返回
	return utils.RDB11.ZCount(c, key, "1", "1").Result()
}

// 判断是否关注该用户
func IsFollow(c *gin.Context, identity string, follower string) (bool, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + identity

	if utils.RDB11.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库
		cnt, err := mysql.IsFollow(follower, identity)
		if err != nil {
			return false, err
		}

		if cnt <= 0 {
			return false, nil
		}

		return true, nil
	}

	// 存在，则获取，并返回
	cnt := utils.RDB11.ZScore(c, key, follower).Val()
	if cnt != 1 {
		return false, nil
	}
	return true, nil
}
