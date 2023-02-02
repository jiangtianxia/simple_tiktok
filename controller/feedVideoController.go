package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"strconv"
	"time"
)

func FeedVideo(c *gin.Context) {

	// 接收参数
	resTime := c.Query("latest_time")
	var latestTime time.Time
	if resTime == "" { // 不填表示当前时间
		latestTime = time.Now()
	} else {
		int64Time, _ := strconv.ParseInt(resTime, 10, 64)
		latestTime = time.Unix(int64Time, 0)
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
			"status_msg":  err.Error(),
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
