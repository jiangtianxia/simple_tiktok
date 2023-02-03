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

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB1中写入hash键值对，但我只设置一个字段
func RedisSetHashRDB1(key string, text string, value string) error {
	pipeline := utils.RDB1.TxPipeline()
	pipeline.HSet(ctx, key, text, value)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB3中写入hash键值对
func RedisSetHashRDB3(key string, maps *map[string]interface{}) error {
	pipeline := utils.RDB3.TxPipeline()
	pipeline.HMSet(ctx, key, *maps)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB4中set键值对添加元素，不用List的原因是防重复
func RedisAddSetRDB4(key string, element string) error {
	pipeline := utils.RDB4.TxPipeline()
	pipeline.SAdd(ctx, key, element)

	_, err := pipeline.Exec(ctx)
	return err
}
