package controller

import (
	"fmt"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"github.com/gin-gonic/gin"
	"simple_tiktok/middlewares"
	"strconv"
)

/* 视频列表
 * @Summary 获取用户发布的视频列表
 * @Tags 基础接口
 */
func GetPublishList(c *gin.Context) {
	userIdS := c.Query("user_id")
	token := c.Query("token")
	// 转换string为uint64
	userIdI, err := strconv.Atoi(userIdS)
	if err != nil {
		logger.SugarLogger.Error("Atoi error:" + err.Error())
		fmt.Println("Atoi Error:" + err.Error())
		GetPublishListReturn(c, -1, fmt.Sprintf("转换字符串 %s 为整型失败", userIdS), nil)
		return
	}
	userId := uint64(userIdI)
	// 验证登陆者token是否可用
	User, err := middlewares.AuthUserCheck(token)
	if err != nil {
		logger.SugarLogger.Error("Get Token Error:" + err.Error())
		fmt.Println("Get Token Error:" + err.Error())
		GetPublishListReturn(c, -1, "Token验证失败", nil)
		return
	}
	// 获取登录者的UserId
	loginUserId := User.Identity
	// 调用service查询结果
	videoList, err := service.GetVideoListByUserId(&userId, &loginUserId)
	if err != nil {
		logger.SugarLogger.Error("Query VideoList Error:" + err.Error())
		fmt.Println("Query VideoList Error:" + err.Error())
		GetPublishListReturn(c, -1, "获取用户视频失败", nil)
		return
	}
	GetPublishListReturn(c, 0, "成功", videoList)
}

// 回传函数
func GetPublishListReturn(c *gin.Context, status_code int, status_msg string, video_list *[]service.Video) {
	// 构造
	var response GetPublishListResponse
	response.StatusCode = status_code
	response.StatusMsg = status_msg
	response.VideoList = *video_list

	c.JSON(200,  response)
}

// 返回参数
type GetPublishListResponse struct {
	StatusCode int `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	VideoList []service.Video `json:"video_list"`
}
