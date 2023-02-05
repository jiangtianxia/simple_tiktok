package redis

import (
	"simple_tiktok/utils"
	"time"
	"github.com/go-redis/redis/v9"
)

/** 
 * 对redis的操作
 * @Author pjh
 * @Summary 
 * @Tags 
 **/

// 往RDB4中set键值对添加元素，不用List的原因是防重复
func RedisAddSetRDB4(key string, element string) error {
	pipeline := utils.RDB4.TxPipeline()
	pipeline.SAdd(ctx, key, element)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB8中把评论id添加到list中
func RedisAddListRBD8(key string, commitId string) error {
	pipeline := utils.RDB8.TxPipeline()
	pipeline.LPush(ctx, key, commitId)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB7中设置评论信息
func RedisSetHashRDB7(key string, value map[string]interface{}) error {
	pipeline := utils.RDB7.TxPipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB5中插入点赞视频的Id和分数
func RedisAddZSetRDB5(key string, value string, status float64) error {
	pipeline := utils.RDB5.TxPipeline()
	pipeline.ZAdd(ctx, key, redis.Z{
		Score: status,
		Member: value,
	})
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB6中把点赞的用户ID添加到list中
func RedisAddListRDB6(key string, value string) error {
	pipeline := utils.RDB6.TxPipeline()
	pipeline.LPush(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}
