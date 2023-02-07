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

	// 2、验证token
	if token == "" {
		UploadResp(c, -1, "无效的Token")
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
