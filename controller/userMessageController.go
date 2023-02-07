package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"simple_tiktok/utils"
)

func SendMessage(c *gin.Context) {
	//token := c.DefaultQuery("token", "")
	tmp, _ := utils.GenerateToken(1, "jack")
	// 验证token
	UserClaims, err := middlewares.AuthUserCheck(tmp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	fromUserId := UserClaims.Identity

	// 接受参数
	toUserId := c.DefaultPostForm("to_user_id", "0")
	actionType := c.DefaultPostForm("action_type", "0")
	content := c.DefaultPostForm("content", "")

	// 把参数传给service层
	err = service.SendMessage(c, fromUserId, toUserId, actionType, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "发送消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "发送消息成功",
	})
}

func MessageRecord(c *gin.Context) {
	//token := c.DefaultQuery("token", "")
	tmp, _ := utils.GenerateToken(1, "jack")
	// 验证token
	UserClaims, err := middlewares.AuthUserCheck(tmp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	fromUserId := UserClaims.Identity

	// 接受参数
	toUserId := c.DefaultQuery("to_user_id", "0")

	// 把参数传给service层
	res, err := service.MessageRecord(c, fromUserId, toUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "查询聊天记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code":  0,
		"status_msg":   "查询聊天记录成功",
		"message_list": res,
	})
}
