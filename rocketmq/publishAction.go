package rocket

import (
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"strconv"
)

func PublishAction(videoInfo models.VideoBasic) {
	videoid := strconv.Itoa(int(videoInfo.Identity))
	userid := strconv.Itoa(int(videoInfo.UserIdentity))
	err := redis.RedisAddVideoList(videoInfo.PublishTime, videoid)
	if err != nil {
		logger.SugarLogger.Error("RedisAddVideoList Error：", err.Error())
	}
	err = redis.RedisAddVideoInfo(videoInfo)
	if err != nil {
		logger.SugarLogger.Error("RedisAddVideoInfo Error：", err.Error())
	}

	err = redis.RedisAddPublishList(userid, videoid)
	if err != nil {
		logger.SugarLogger.Error("RedisAddPublishList Error：", err.Error())
	}

	err = redis.RedisAddFavoriteUser(videoid, userid, 0)
	if err != nil {
		logger.SugarLogger.Error("RedisAddFavoriteUser Error：", err.Error())
	}
}
