package redis

import (
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

func RedisAddUserMessageHash(ctx *gin.Context, list []models.UserMessage) error {
	hashKey := viper.GetString("redis.KeyUserMessageHashPrefix")

	var key string
	var value map[string]interface{}
	// 开启事务
	pipeline := utils.RDB13.TxPipeline()
	for i := 0; i < len(list); i++ {
		key = hashKey + strconv.FormatUint(list[i].Identity, 10)
		value = map[string]interface{}{
			"id":          list[i].Identity,
			"content":     list[i].Content,
			"create_time": list[i].CreateTime,
		}
		pipeline.HSet(ctx, key, value)
		pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func RedisAddUserMessageSet(ctx *gin.Context, key string, list []models.UserMessage) error {
	// 开启事务
	pipeline := utils.RDB12.TxPipeline()

	for i := 0; i < len(list); i++ {
		// // createTime, err := time.Parse("2006-01-02 15:04:05", list[i].CreateTime)
		// if err != nil {
		// 	logger.SugarLogger.Error(err)
		// 	return err
		// }
		pipeline.ZAdd(ctx, key, redis.Z{
			Score:  float64(list[i].CreateTime),
			Member: strconv.FormatUint(list[i].Identity, 10),
		})
		pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
	}

	_, err := pipeline.Exec(ctx)
	return err
}
