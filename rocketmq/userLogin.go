package rocket

import (
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
)

func UserLogin(userLogin models.UserBasic) {
	if err := redis.RedisLogin(userLogin); err != nil {
		logger.SugarLogger.Error("RedisLogin Error: ", err.Error())
	}
}
