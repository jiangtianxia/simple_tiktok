package controller

import (
	"encoding/json"
	"net/http"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 返回结构体
type SendMessageRespStruct struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// 传入参数返回
func SendMessageResp(c *gin.Context, code int, msg string) {
	h := &SendMessageRespStruct{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, h)
}

// 发送消息接收参数结构体
type SendMessageReqStruct struct {
	FromUserId uint64
	ToUserId   string
	ActionType string
	Content    string
}

// SendMessage
// @Summary 发送消息
// @Tags 社交接口
// @Param token query string true "token"
// @Param to_user_id query string true "用户id"
// @Param action_type query string true "1-发送消息"
// @Param content query string true "消息内容"
// @Success 200 {object} SendMessageRespStruct
// @Router /message/action/ [post]
func SendMessage(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	// 验证token
	UserClaims, err := middlewares.AuthUserCheck(token)
	if err != nil {
		SendMessageResp(c, -1, "无效token")
		return
	}
	fromUserId := UserClaims.Identity

	// 接受参数
	toUserId := c.DefaultQuery("to_user_id", "0")
	actionType := c.DefaultQuery("action_type", "0")
	content := c.DefaultQuery("content", "")

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
		logger.SugarLogger.Error("发送失败， error:", err.Error())
		SendMessageResp(c, -1, "发送失败")
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
				SendMessageResp(c, code, info["status_msg"])
				return
			}

			time.Sleep(time.Second / 15)
		}

		// 20秒后，还是没有结果，则返回请求超时
		SendMessageResp(c, -1, "请求超时")
		return
	}

	SendMessageResp(c, -1, "发送失败")
}

// 返回结构体
type MessageRecordRespStruct struct {
	Code        int               `json:"status_code"`
	Msg         string            `json:"status_msg"`
	MessageList []service.Message `json:"message_list"`
}

// 传入参数返回
func MessageRecordResp(c *gin.Context, code int, msg string, messageList []service.Message) {
	h := &MessageRecordRespStruct{
		Code:        code,
		Msg:         msg,
		MessageList: messageList,
	}

	c.JSON(http.StatusOK, h)
}

// MessageRecord
// @Summary 聊天记录
// @Tags 社交接口
// @Param token query string true "token"
// @Param to_user_id query string true "用户id"
// @Success 200 {object} MessageRecordRespStruct
// @Router /message/chat/ [get]
func MessageRecord(c *gin.Context) {
	token := c.DefaultQuery("token", "")

	// 验证token
	UserClaims, err := middlewares.AuthUserCheck(token)
	if err != nil {
		MessageRecordResp(c, -1, "无效token", []service.Message{})
		return
	}
	fromUserId := UserClaims.Identity

	// 接受参数
	toUserId := c.DefaultQuery("to_user_id", "0")

	// 把参数传给service层
	res, err := service.MessageRecord(c, fromUserId, toUserId)
	if err != nil {
		MessageRecordResp(c, -1, "查询聊天记录失败", []service.Message{})
		return
	}
	MessageRecordResp(c, 0, "查询聊天记录成功", res)
}
