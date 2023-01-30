package controller

import (
	"fmt"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"simple_tiktok/dao/mysql"
	"github.com/gin-gonic/gin"
	"simple_tiktok/middlewares"
)

// 返回参数
type GetPublishListResponse struct {
	StatusCode int `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	VideoList []mysql.Video `json:"video_list"`
}

/* 视频列表
 * @Summary 获取用户发布的视频列表
 * @Tags 基础接口
 */
func GetPublishList(c *gin.Context) {
	userId := c.Query("user_id")
	token := c.Query("token")
	// 验证token是否可用
	_, err := middlewares.AuthUserCheck(token)
	if err != nil {
		logger.SugarLogger.Error("Get Token Error:" + err.Error())
		fmt.Println("Get Token Error:" + err.Error())
		GetPublishListReturn(c, -1, "Token验证失败", nil)
		return
	}
	// 调用service查询
	videoList, err := service.GetVideoListByUserId(&userId)
	if err != nil {
		logger.SugarLogger.Error("Query VideoList Error:" + err.Error())
		fmt.Println("Query VideoList Error:" + err.Error())
		GetPublishListReturn(c, -1, "获取用户视频失败", nil)
		return
	}
	GetPublishListReturn(c, 0, "成功", videoList)
}

// 回传函数
func GetPublishListReturn(c *gin.Context, status_code int, status_msg string, video_list *[]mysql.Video) {
	var response GetPublishListResponse
	response.StatusCode = status_code
	response.StatusMsg = status_msg
	response.VideoList = *video_list

	c.JSON(200,  response)
}
