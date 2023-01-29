package controller

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"simple_tiktok/logger"
	"simple_tiktok/models"
)

// 处理函数
func GetPublishList(c *gin.Context) {
	var query models.GetPublishListQuery
	// 绑定query参数
	err := c.ShouldBind(&query)
	if err != nil {
		logger.SugarLogger.Error("Get FormFile Error:" + err.Error())
		fmt.Println("Get FormFile Error:" + err.Error())
		GetPublishListReturn(c, -1, "获取参数失败", nil)
		return
	}
	// token或者userid为空
	if query.Token == "" || query.UserId == "" {
		GetPublishListReturn(c, -1, "请求参数错误", nil)
		return
	}
	// 调用service查询
	// -TODO
}

// 回传函数
func GetPublishListReturn(c *gin.Context, status_code int, status_msg string, video_list *[]models.Video) {
	var response models.GetPublishListResponse
	response.StatusCode = status_code
	response.StatusMsg = status_msg
	response.VideoList = *video_list

	c.JSON(200,  response)
}
