package middlewares

import (
	"net/http"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 熔断器的使用
 * @Date 22:00 2023/1/31
 **/
func GinCircuitBreaker(c *gin.Context) {
	url := c.FullPath()
	if !utils.CB.CanAccess(url) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"msg": "服务器错误，请稍后重试",
		})
		return
	}

	c.Next()

	if c.Writer.Status() > 410 {
		utils.CB.Failed(url)
		return
	}

	utils.CB.Succeed(url)
}
