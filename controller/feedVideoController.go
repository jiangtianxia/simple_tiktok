package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"strconv"
	"time"
)

// FeedVideo
// @Summary 视频流
// @Tags 基础接口
// @Param latest_time query string false "latest_time"
// @Param token query string false "token"
// @Success 200 {string} status_code status_msg
// @Router /feed [get]
func FeedVideo(c *gin.Context) {

	// 接收参数
	resTime := c.Query("latest_time")
	var latestTime int64
	if resTime == "" { // 不填表示当前时间
		latestTime = time.Now().Unix()
	} else {
		latestTime, _ = strconv.ParseInt(resTime, 10, 64)
	}
	token := c.Query("token")

	if token != "" {
		// 验证token
		_, err := middlewares.AuthUserCheck(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": -1,
				"status_msg":  err.Error(),
			})
			return
		}
	}

	res, nextTime, err := service.FeedVideo(c, latestTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "获取视频失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"next_time":   nextTime,
		"video_list":  res,
	})
}
