package service

import (
	"errors"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"time"
)

func DealFavourite(token string, video_id string, action_type string) error {
	userClaim, err := middlewares.AuthUserCheck(token)
	if err != nil {
		logger.SugarLogger.Error("token error:" + err.Error())
	}
	exist := mysql.IsFavourite(userClaim.Identity, video_id)
	//key := viper.GetString("redis.KetFavoriteSetPrefix") + video_id
	if exist && action_type == "2" {
		//删除缓存
		err := redis.RedisDeleteFavoriteUser(video_id, userClaim.Identity)
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}
		//修改数据库
		err = mysql.UpdateFavourite(video_id, userClaim.Identity, "0")
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}

		//休眠2s，再次删除缓存
		time.Sleep(time.Second * 2)
		err = redis.RedisDeleteFavoriteUser(video_id, userClaim.Identity)
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}

	} else if !exist && action_type == "1" {
		//删除缓存
		err := redis.RedisDeleteFavoriteUser(video_id, userClaim.Identity)
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}
		//修改数据库
		err = mysql.UpdateFavourite(video_id, userClaim.Identity, "1")
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}

		//休眠2s，再次删除缓存
		time.Sleep(time.Second * 2)
		err = redis.RedisDeleteFavoriteUser(video_id, userClaim.Identity)
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			return err
		}

	} else {
		logger.SugarLogger.Error("参数 error")
		return errors.New("参数错误")
	}

	return nil
}
