package rocket

import (
	"fmt"
	"github.com/spf13/viper"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"strconv"
)

func UserInfoAction(userinfo models.UserBasic) {
	var newCathe = map[string]interface{}{
		"identity": userinfo.Identity,
		"username": userinfo.Username,
	}
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + strconv.Itoa(int(userinfo.Identity))
	err := myRedis.RedisAddUserInfo(hashKey, newCathe)
	if err != nil {
		logger.SugarLogger.Error(err)
	}
	fmt.Println("消息队列处理成功")
}
