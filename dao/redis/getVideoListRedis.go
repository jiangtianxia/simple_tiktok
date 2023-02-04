package redis

import (
	"simple_tiktok/utils"
	"time"
)

/** 
 * 对redis的操作
 * @Author pjh
 * @Summary 
 * @Tags 
 **/

// 往RDB0中写入string键值对
func RedisAddStringRDB0(key string, value string) error {
	pipeline := utils.RDB0.TxPipeline()
	pipeline.Set(ctx, key, value, 10 * time.Hour)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

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
