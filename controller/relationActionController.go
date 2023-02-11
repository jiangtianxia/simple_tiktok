package controller

import (
	"encoding/json"
	"net/http"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 关注操作接口
 * @Date 15:00 2023/2/11
 **/
// 返回结构体
type FollowRespStruct struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// 传入参数返回
func FollowResp(c *gin.Context, code int, msg string) {
	h := &FollowRespStruct{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, h)
}

// 接收参数结构体
type FollowReqStruct struct {
	UserId     string
	ToUserId   string
	ActionType int
}

// Follow
// @Summary 关注操作
// @Tags 社交接口
// @Param token query string true "token"
// @Param to_user_id query string true "用户id"
// @Param action_type query string true "关注操作"
// @Success 200 {object} FollowRespStruct
// @Router /relation/action/ [post]
func Follow(c *gin.Context) {
	// 1、获取参数
	token := c.DefaultQuery("token", "")
	to_user_id := c.DefaultQuery("to_user_id", "")
	action_type := c.DefaultQuery("action_type", "")

	// 2、验证参数
	if token == "" || to_user_id == "" || action_type == "" {
		FollowResp(c, -1, "请求参数错误")
		return
	}

	actionType, _ := strconv.Atoi(action_type)
	if actionType != 1 && actionType != 2 {
		FollowResp(c, -1, "请求参数错误")
		return
	}

	user_id, _ := strconv.Atoi(to_user_id)
	cnt, err := mysql.FindUserByIdentityCount(uint64(user_id))
	if err != nil {
		logger.SugarLogger.Error("FindUserByIdentity1 Error：", err.Error())
		FollowResp(c, -1, "请求参数错误")
		return
	}
	if cnt == 0 {
		FollowResp(c, -1, "用户不存在")
		return
	}

	// 3、验证token
	t, _ := utils.GenerateToken(1, "test")
	UserClaims, err := middlewares.AuthUserCheck(t)
	if err != nil {
		FollowResp(c, -1, "无效token")
		return
	}

	// 4、将消息发送到消息队列
	Info := FollowReqStruct{
		UserId:     strconv.Itoa(int(UserClaims.Identity)),
		ToUserId:   to_user_id,
		ActionType: actionType,
	}
	data, _ := json.Marshal(Info)

	producer := viper.GetString("rocketmq.serverProducer")
	topic := viper.GetString("rocketmq.ServerTopic")
	tag := viper.GetString("rocketmq.serverFollowTag")
	// 发送消息
	res, err := utils.SendMsg(c, producer, topic, tag, data)
	if err != nil {
		logger.SugarLogger.Error("发送消息失败， error:", err.Error())
		FollowResp(c, -1, "操作失败")
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
				FollowResp(c, code, info["status_msg"])
				return
			}

			time.Sleep(time.Second / 15)
		}

		// 20秒后，还是没有结果，则返回请求超时
		FollowResp(c, -1, "请求超时！！！")
		return
	}

	FollowResp(c, -1, "操作失败")
}
