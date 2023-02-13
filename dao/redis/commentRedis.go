package redis

import (
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 往RDB8中把评论id添加到list中
func RedisAddListRBD8(ctx *gin.Context, key string, commitId string) error {
	pipeline := utils.RDB8.TxPipeline()
	pipeline.LPush(ctx, key, commitId)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}

// 往RDB7中设置评论信息
func RedisSetHashRDB7(ctx *gin.Context, key string, value map[string]interface{}) error {
	pipeline := utils.RDB7.TxPipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Hour*24*5)

	_, err := pipeline.Exec(ctx)
	return err
}
