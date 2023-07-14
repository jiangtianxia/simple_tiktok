package controller

import (
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 返回结构体
type FeedVideoRespStruct struct {
	Code      int                 `json:"status_code"`
	Msg       string              `json:"status_msg"`
	NextTime  int64               `json:"next_time"`
	VideoList []service.VideoInfo `json:"video_list"`
}

// 传入参数返回
func FeedVideoResp(c *gin.Context, code int, msg string, nextTime int64, videoList []service.VideoInfo) {
	h := &FeedVideoRespStruct{
		Code:      code,
		Msg:       msg,
		NextTime:  nextTime,
		VideoList: videoList,
	}

	c.JSON(http.StatusOK, h)
}

// FeedVideo
//	@Summary	视频流接口
//	@Tags		基础接口
//	@Param		latest_time	query		string	false	"latest_time"
//	@Param		token		query		string	false	"token"
//	@Success	200			{object}	FeedVideoRespStruct
//	@Router		/feed [get]
func FeedVideo(c *gin.Context) {
	// 接收参数
	resTime := c.DefaultQuery("latest_time", "")

	var latestTime int64
	if resTime == "" { // 不填表示当前时间
		latestTime = time.Now().Unix()
	} else {
		latestTime, _ = strconv.ParseInt(resTime, 10, 64)
	}
	token := c.DefaultQuery("token", "")
	var userId uint64
	userId = 0

	if token != "" {
		// 验证token
		userClaim, err := middlewares.AuthUserCheck(token)
		userId = userClaim.Identity
		if err != nil {
			FeedVideoResp(c, -1, err.Error(), 0, []service.VideoInfo{})
			return
		}
	}

	res, nextTime, err := service.FeedVideo(c, userId, latestTime)
	if err != nil {
		FeedVideoResp(c, -1, "获取视频失败", 0, []service.VideoInfo{})
		return
	}

	FeedVideoResp(c, 0, "获取视频成功", nextTime, res)
}
