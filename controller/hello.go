package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello
//	@Tags		公共接口
//	@Summary	首页
//	@Success	200	{string}	hello	world
//	@Router		/hello [get]
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}
