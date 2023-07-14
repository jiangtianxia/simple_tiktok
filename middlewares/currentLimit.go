package middlewares

import (
	"net/http"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 令牌桶限流策略
 * @Date 22:00 2023/1/31
 **/
func CurrentLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 判断令牌桶中是否有令牌
		if !utils.Bucket.Allow() {
			// 不允许访问
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "服务器繁忙，请稍后重试",
			})
			return
		}
		ctx.Next()
	}
}
