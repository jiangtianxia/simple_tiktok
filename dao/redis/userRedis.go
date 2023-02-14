package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author Xiaoyu Zhang
 * @Description 新增用户信息缓存
 * @Date 14:00 2023/1/31
 **/
func RedisAddUserInfo(c *gin.Context, key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB1.TxPipeline()
	pipeline.HSet(c, key, value)
	pipeline.Expire(c, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(c)
	return err
}

/**
 * @Author jiang
 * @Description 新增用户信息缓存
 * @Date 13:00 2023/2/12
 **/
func RedisAddUserInfoHash(c *gin.Context, key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB1.TxPipeline()
	pipeline.HSet(c, key, value)
	pipeline.Expire(c, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(c)
	return err
}
