package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/service"
)

func Favourite(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	// 验证token
	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "无效的Token",
		})
		return
	}

	if action_type != "2" && action_type != "1" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "无效的action_type",
		})
		return
	}
	err := service.DealFavourite(token, video_id, action_type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "2",
	})

}
