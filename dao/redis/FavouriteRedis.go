package redis

import (
	"simple_tiktok/utils"
	"strconv"
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
	pipeline.HDel(ctx, key, strconv.Itoa(int(userId)))
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
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB6中把点赞的用户ID添加到list中
func RedisAddListRDB6(ctx *gin.Context, key string, value string) error {
	pipeline := utils.RDB6.TxPipeline()
	pipeline.LPush(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}

/**
 * @Author jiang
 * @Description 使用有序集合存储视频的点赞用户sorted set
 * @Date 19:00 2023/1/29
 **/
func RedisAddFavoriteUser(ctx *gin.Context, videoId string, userId string, score int) error {
	// 1、获取前缀，拼接key
	key := viper.GetString("redis.KetFavoriteSetPrefix") + videoId

	// 2、创建成员并添加数据
	//  ZADD KEY_NAME SCORE1 VALUE1.. SCOREN VALUEN
	// 开启事务
	pipeline := utils.RDB5.Pipeline()
	pipeline.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: userId,
	})
	pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)

	_, err := pipeline.Exec(ctx)
	return err
}
