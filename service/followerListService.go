package service

import (
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 粉丝列表
 * @Date 17:00 2023/2/12
 **/
func FollowerListService(c *gin.Context, userId uint64) ([]Author, error) {
	data := make([]Author, 0)

	// 1、判断缓存中是否存在数据
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + strconv.Itoa(int(userId))

	if utils.RDB11.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库
		// 1）根据id，查询关注的用户信息
		followerList, err := mysql.FindFollower(strconv.Itoa(int(userId)))
		if err != nil {
			logger.SugarLogger.Error("FindFollower Error：", err.Error())
			return nil, err
		}

		// 2）循环遍历list，获取用户详细信息
		pipeline := utils.RDB11.Pipeline()
		for _, follower := range followerList {
			// 获取用户名
			username, err := GetUsername(c, strconv.Itoa(int(follower.UserIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetUsername Error：", err.Error())
				return nil, err
			}

			// 获取关注数
			followCount, err := GetFollowCount(c, strconv.Itoa(int(follower.UserIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
				return nil, err
			}

			// 获取粉丝数
			followerCount, err := GetFollowerCount(c, strconv.Itoa(int(follower.UserIdentity)))
			if err != nil {
				logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
				return nil, err
			}

			// 判断是否关注该粉丝
			isFollow, err := IsFollow(c, strconv.Itoa(int(follower.UserIdentity)), strconv.Itoa(int(userId)))
			if err != nil {
				logger.SugarLogger.Error("IsFollow Error：", err.Error())
				return nil, err
			}

			totalFavourited, workCount, FavouriteCount, err := GetTotalFavouritedANDWorkCountANDFavoriteCount(follower.UserIdentity)
			if err != nil {
				logger.SugarLogger.Error("GetTotalFavouritedANDWorkCountANDFavoriteCount Error：", err.Error())
				return nil, err
			}

			userInfo := Author{
				Id:              follower.UserIdentity,
				Name:            username,
				FollowCount:     followCount,
				FollowerCount:   followerCount,
				IsFollow:        isFollow,
				Avatar:          viper.GetString("defaultAvatarUrl"),
				BackGroundImage: viper.GetString("defaultBackGroudImage"),
				Signature:       viper.GetString("defaultSignature"),
				TotalFavorited:  totalFavourited,
				WorkCount:       workCount,
				FavoriteCount:   FavouriteCount,
			}
			pipeline.ZAdd(c, key, redis.Z{
				Member: follower.UserIdentity,
				Score:  1,
			})
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

	// 存在，则获取缓存中的粉丝数据
	// ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
	followerList, err := utils.RDB11.ZRangeByScore(c, key, &redis.ZRangeBy{
		Min: "1",
		Max: "1",
	}).Result()
	if err != nil {
		logger.SugarLogger.Error("ZRangeByScore Error：", err.Error())
		return nil, err
	}

	for _, followerIdentity := range followerList {
		// 获取用户名
		username, err := GetUsername(c, followerIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetUsername Error：", err.Error())
			return nil, err
		}

		// 获取关注数
		followCount, err := GetFollowCount(c, followerIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
			return nil, err
		}

		// 获取粉丝数
		followerCount, err := GetFollowerCount(c, followerIdentity)
		if err != nil {
			logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
			return nil, err
		}

		// 判断是否关注该粉丝
		isFollow, err := IsFollow(c, followerIdentity, strconv.Itoa(int(userId)))
		if err != nil {
			logger.SugarLogger.Error("IsFollow Error：", err.Error())
			return nil, err
		}

		userIdentity64, _ := strconv.Atoi(followerIdentity)
		totalFavourited, workCount, FavouriteCount, err := GetTotalFavouritedANDWorkCountANDFavoriteCount(uint64(userIdentity64))
		if err != nil {
			logger.SugarLogger.Error("GetTotalFavouritedANDWorkCountANDFavoriteCount Error：", err.Error())
			return nil, err
		}

		id, _ := strconv.Atoi(followerIdentity)
		userInfo := Author{
			Id:              uint64(id),
			Name:            username,
			FollowCount:     followCount,
			FollowerCount:   followerCount,
			IsFollow:        isFollow,
			Avatar:          viper.GetString("defaultAvatarUrl"),
			BackGroundImage: viper.GetString("defaultBackGroudImage"),
			Signature:       viper.GetString("defaultSignature"),
			TotalFavorited:  totalFavourited,
			WorkCount:       workCount,
			FavoriteCount:   FavouriteCount,
		}
		data = append(data, userInfo)
	}
	return data, nil
}
