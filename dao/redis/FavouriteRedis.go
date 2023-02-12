package redis

import (
	"github.com/spf13/viper"
	"simple_tiktok/utils"
	"time"
)

func RedisDeleteFavoriteUser(videoId string, userId uint64) error {
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
