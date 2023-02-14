package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 往RDB8中把评论id添加到list中
func RedisAddListRBD8(ctx *gin.Context, key string, commitId string) error {
	pipeline := utils.RDB8.TxPipeline()
	pipeline.LPush(ctx, key, commitId)
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB7中设置评论信息
func RedisSetHashRDB7(ctx *gin.Context, key string, value map[string]interface{}) error {
	pipeline := utils.RDB7.TxPipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}

/**
 * Creator: lx
 * Last Editor: lx
 * Description: 新增评论信息缓存
 **/

func RedisAddCommentInfo(ctx *gin.Context, key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB7.TxPipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}
