package redis

import (
	"github.com/gin-gonic/gin"
	"simple_tiktok/utils"
	"time"
)

func RedisLogin(c *gin.Context, key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB.TxPipeline()
	pipeline.HSet(c, key, value)
	pipeline.Expire(c, key, time.Hour*24*5)

	_, err := pipeline.Exec(c)
	return err

}
