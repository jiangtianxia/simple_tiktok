package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/spf13/viper"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: 新增评论信息缓存
 **/

 func RedisAddCommentInfo(key string, value map[string]interface{}) error {
	// 开启事务
	pipeline := utils.RDB7.TxPipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))

	_, err := pipeline.Exec(ctx)
	return err
}
