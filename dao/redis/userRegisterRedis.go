package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: 用户注册缓存
 **/

func RedisUserRegister(c *gin.Context, key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB1.TxPipeline()
	pipeline.HSet(c, key, value)
	pipeline.Expire(c, key, time.Hour*24*5)

	_, err := pipeline.Exec(c)
	return err
}
