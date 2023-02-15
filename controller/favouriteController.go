package controller

import (
	"encoding/json"
	"net/http"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 返回结构体
type FavouriteRespStruct struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// 传入参数返回
func FavouriteResp(c *gin.Context, code int, msg string) {
	h := &FavouriteRespStruct{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, h)
}

// 接收参数结构体
type FavouriteReqStruct struct {
	UserId     uint64
	VideoId    string
	ActionType string
}

// Favourite
// @Summary 赞操作
// @Tags 互动接口
// @Param token query string true "token"
// @Param video_id query string true "视频id"
// @Param action_type query string true "赞操作"
// @Success 200 {object} FavouriteRespStruct
// @Router /favorite/action/ [post]
func Favourite(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	// 验证参数
	if action_type != "2" && action_type != "1" {
		FavouriteResp(c, -1, "参数不正确")
		return
	}

	if token == "" || video_id == "" {
		FavouriteResp(c, -1, "参数不正确")
		return
	}

	// 验证token
	UserClaims, err := middlewares.AuthUserCheck(token)
	if err != nil {
		FavouriteResp(c, -1, "无效token")
		return
	}

	// 将消息发送到消息队列
	Info := FavouriteReqStruct{
		UserId:     UserClaims.Identity,
		VideoId:    video_id,
		ActionType: action_type,
	}
	data, _ := json.Marshal(Info)

	producer := viper.GetString("rocketmq.serverProducer")
	topic := viper.GetString("rocketmq.ServerTopic")
	tag := viper.GetString("rocketmq.serverFavouriteTag")
	// 发送消息
	res, err := utils.SendMsg(c, producer, topic, tag, data)
	if err != nil {
		logger.SugarLogger.Error("发送消息失败， error:", err.Error())
		FavouriteResp(c, -1, "操作失败")
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
				FavouriteResp(c, code, info["status_msg"])
				return
			}

			time.Sleep(time.Second / 15)
		}

		// 20秒后，还是没有结果，则返回请求超时
		FavouriteResp(c, -1, "请求超时！！！")
		return
	}

	FavouriteResp(c, -1, "操作失败")
}
