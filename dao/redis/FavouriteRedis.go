package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

func RedisDeleteFavoriteUser(ctx *gin.Context, videoId string, userId uint64) error {
	// 1、获取前缀，拼接key
	key := viper.GetString("redis.KetFavoriteSetPrefix") + videoId

	// 2、删除成员及数据
	//  ZADD KEY_NAME SCORE1 VALUE1.. SCOREN VALUEN
	// 开启事务
	pipeline := utils.RDB5.Pipeline()
	pipeline.HDel(ctx, key, string(userId))
	pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB5中插入点赞视频的Id和分数
func RedisAddZSetRDB5(ctx *gin.Context, key string, value string, status float64) error {
	pipeline := utils.RDB5.TxPipeline()
	pipeline.ZAdd(ctx, key, redis.Z{
		Score:  status,
		Member: value,
	})
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB6中把点赞的用户ID添加到list中
func RedisAddListRDB6(ctx *gin.Context, key string, value string) error {
	pipeline := utils.RDB6.TxPipeline()
	pipeline.LPush(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}
