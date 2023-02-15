package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"strconv"
)

// 返回结构体
type CommentListStruct struct {
	Code        int                   `json:"status_code"`
	Msg         string                `json:"status_msg"`
	CommentList []service.CommentInfo `json:"comment_list"`
}

// 传入参数返回
func CommentListResp(c *gin.Context, code int, msg string, commentList []service.CommentInfo) {
	h := &CommentListStruct{
		Code:        code,
		Msg:         msg,
		CommentList: commentList,
	}

	c.JSON(http.StatusOK, h)
}

func CommentList(c *gin.Context) {

	token := c.DefaultQuery("token", "")
	videoIdStr := c.Query("video_id")

	fmt.Print(token)

	//没有token 查看不了评论
	if token == "" {
		CommentListResp(c, -1, "获取评论失败", []service.CommentInfo{})
		return
	}

	// 验证token
	userClaim, err := middlewares.AuthUserCheck(token)
	userId := userClaim.Identity
	if err != nil {
		CommentListResp(c, -1, "获取评论失败", []service.CommentInfo{})
		return
	}
	intNum, _ := strconv.Atoi(videoIdStr)
	videoId := uint64(intNum)

	res, err := service.CommentList(c, userId, videoId)
	if err != nil {
		CommentListResp(c, -1, "获取评论失败", []service.CommentInfo{})
		return
	}
	CommentListResp(c, 0, "获取评论成功", res)
}
