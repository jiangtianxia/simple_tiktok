package redis

import (
	"context"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

func RedisAddUserMessageHash(list []models.UserMessage) error {
	hashKey := viper.GetString("redis.KeyUserMessageHashPrefix")

	var key string
	var value map[string]interface{}
	// 开启事务
	pipeline := utils.RDB13.TxPipeline()
	var c = context.Background()
	for i := 0; i < len(list); i++ {
		key = hashKey + strconv.FormatUint(list[i].Identity, 10)
		value = map[string]interface{}{
			"id":          list[i].Identity,
			"content":     list[i].Content,
			"create_time": list[i].CreateTime,
		}
		pipeline.HSet(c, key, value)
		pipeline.Expire(c, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
	}

	_, err := pipeline.Exec(c)
	return err
}

func RedisAddUserMessageSet(key string, list []models.UserMessage) error {
	// 开启事务
	pipeline := utils.RDB12.TxPipeline()
	var c = context.Background()
	for i := 0; i < len(list); i++ {
		// // createTime, err := time.Parse("2006-01-02 15:04:05", list[i].CreateTime)
		// if err != nil {
		// 	logger.SugarLogger.Error(err)
		// 	return err
		// }
		pipeline.ZAdd(c, key, redis.Z{
			Score:  float64(list[i].CreateTime),
			Member: strconv.FormatUint(list[i].Identity, 10),
		})
		pipeline.Expire(c, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
	}

	_, err := pipeline.Exec(c)
	return err
}
