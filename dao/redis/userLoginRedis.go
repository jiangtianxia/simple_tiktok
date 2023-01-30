package redis

import (
	"github.com/spf13/viper"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"
	"time"
)

//func RedisLogin(c *gin.Context, key string, value map[string]interface{}) error {
//	// 开启事务
//	pipeline := utils.RDB.TxPipeline()
//	pipeline.HSet(c, key, value)
//	pipeline.Expire(c, key, time.Hour*24*5)
//
//	_, err := pipeline.Exec(c)
//	return err
//
//}

/**
 * @Auther duan
 * @Description 存储用户登录信息到hash集合
 * @Date 20:00 2023/1/30
 **/
func RedisLogin(userLogin models.UserBasic) error {
	key := viper.GetString("redis.KeyVideoInfoHashPrefix") + strconv.Itoa(int(userLogin.Identity))
	value := map[string]interface{}{
		"id":       userLogin.Identity,
		"username": userLogin.Username,
		"password": userLogin.Password,
	}
	// 开启事务
	pipeline := utils.RDB3.Pipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)

	_, err := pipeline.Exec(ctx)
	return err

}
