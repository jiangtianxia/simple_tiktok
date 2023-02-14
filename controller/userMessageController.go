package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"simple_tiktok/utils"
	"strconv"
	"time"
)

// 发送消息接收参数结构体
type SendMessageReqStruct struct {
	FromUserId uint64
	ToUserId   string
	ActionType string
	Content    string
}

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

	info := SendMessageReqStruct{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		ActionType: actionType,
		Content:    content,
	}
	data, _ := json.Marshal(info)

	producer := viper.GetString("rocketmq.serverProducer")
	topic := viper.GetString("rocketmq.ServerTopic")
	tag := viper.GetString("rocketmq.serverSendMessageTags")
	// 发送消息
	res, err := utils.SendMsg(c, producer, topic, tag, data)
	if err != nil {
		logger.SugarLogger.Error("发送消息失败， error:", err.Error())
		return
	}

	if res.Status == 0 {
		for i := 0; i < 300; i++ {
			// hash msgid req
			// 根据msgid，查询redis缓存中是否存在数据，如果存在则将结果返回
			key := res.MsgID
			// 判断当前key，是否存在
			if utils.RDB0.Exists(c, key).Val() == 1 {
				// 存在，则获取结果返回
				info, _ := utils.RDB0.HGetAll(c, key).Result()
				code, _ := strconv.Atoi(info["status_code"])

				c.JSON(http.StatusOK, gin.H{
					"status_code": code,
					"status_msg":  info["status_msg"],
				})
				return
			}

			time.Sleep(time.Second / 15)
		}

		// 20秒后，还是没有结果，则返回请求超时
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "请求超时",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": -1,
		"status_msg":  "发送信息失败",
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
