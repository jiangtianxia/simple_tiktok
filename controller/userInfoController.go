package controller

import (
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

/**
 * @Author Xiaoyu Zhang
 * @Description 用户信息接口
 * @Date 14:00 2023/1/31
 **/
func UserInfo(c *gin.Context) {
	tmp, _ := utils.GenerateToken(1, "merry")
	// 接受参数
	userId := c.DefaultQuery("user_id", "0")
	//token := c.DefaultQuery("token", "")

	// 验证token
	_, err := middlewares.AuthUserCheck(tmp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}

	// 布隆过滤器
	//filterRes, _ := utils.BloomFilterCheck(userId)
	//if !filterRes {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status_code": -1,
	//		"status_msg":  "no such user",
	//	})
	//	return
	//}

	// 把参数传给service层
	res, err := service.UserInfo(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "string",
		"user":        res,
	})
}
