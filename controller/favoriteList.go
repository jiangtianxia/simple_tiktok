package controller

import (
	"simple_tiktok/service"
	"simple_tiktok/middlewares"
	"simple_tiktok/logger"
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
)

// @Summary 喜欢视频列表接口
// @Produce json
// @Param token query string true "登录用户的token"
// @Param user_id query string true "查找目标用户的id"
// @Success 200 {object} controller.GetPublishListResponse "status_msg为成功"
// @Failue 200 {object} controller.GetPublishListResponse "status_msg包含失败原因"
// @Router /favorite/list/ [get]
func GetFavoriteList(c *gin.Context) {
	userIdS := c.Query("user_id")
	token := c.Query("token")
	// 转换string为uint64
	userIdI, err := strconv.Atoi(userIdS)
	if err != nil {
		logger.SugarLogger.Error("Atoi error:" + err.Error())
		fmt.Println("Atoi Error:" + err.Error())
		GetPublishListReturn(c, -1, "用户不存在", nil)
		return
	}
	userId := uint64(userIdI)
	// 验证登陆者token是否可用
	User, err := middlewares.AuthUserCheck(token)
	if err != nil {
		logger.SugarLogger.Error("Get Token Error:" + err.Error())
		fmt.Println("Get Token Error:" + err.Error())
		GetPublishListReturn(c, -1, "Token失效", nil)
		return
	}
	// 获取登录者的UserId
	loginUserId := User.Identity
	// 调用service查询结果
	videoList, err := service.GetFavoriteListByUserId(c, &userId, &loginUserId)
	if err != nil {
		logger.SugarLogger.Error("Query VideoList Error:" + err.Error())
		fmt.Println("Query VideoList Error:" + err.Error())
		GetPublishListReturn(c, -1, "获取用户视频失败", nil)
		return
	}
	GetPublishListReturn(c, 0, "成功", videoList)
}
