package controller

import (
	"encoding/json"
	"net/http"

	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/models"
	"simple_tiktok/service"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: controller层，解析参数，处理参数，并打包传给service层
 **/

// 返回体
type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    models.CommentVideo `json:"comment"`
}

// 发表，删除评论 /comment/action
// CommentAction
// @Summary 评论操作
// @Tags 互动接口
// @Param token query string true "token"
// @Param video_id query string true "视频id"
// @Param action_type query string true "评论操作"
// @Param comment_text query string false "评论内容"
// @Param comment_id query string false "评论id"
// @Router /comment/action/ [post]
func CommentAction(c *gin.Context) {
	//参数获取
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	comment_text := c.Query("comment_text")
	comment_id := c.Query("comment_id")

	// 参数处理
	// 验证用户token
	user, err := middlewares.AuthUserCheck(token)
	if user == nil || user.Identity == 0 || user.Issuer != "simple_tiktok" || user.Username == "" || err != nil {
		logger.SugarLogger.Error("Unauthorized User")
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户错误",
		})
		return
	}

	videoidentity, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "评论id格式错误",
		})
		return
	}
	commentidentity, err := strconv.Atoi(comment_id)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频id格式错误",
		})
		return
	}
	actiontype, err := strconv.Atoi(action_type)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "操作类型错误",
		})
		return
	}
	// 处理时间格式
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	timeNow := month + "-" + day

	// 数据打包，将消息发送到消息队列
	pack := models.CommentVideo{
		VideoIdentity: (uint64)(videoidentity),
		UserIdentity:  (uint64)(commentidentity),
		Text:          comment_text,
		CommentTime:   timeNow,
	}
	req := service.CommentActionRequire{
		Model:      pack,
		ActionType: actiontype,
	}
	data, _ := json.Marshal(req)

	producer := viper.GetString("rocketmq.serverProducer")
	topic := viper.GetString("rocketmq.ServerTopic")
	tag := viper.GetString("rocketmq.serverSendCommentTag")
	// 发送消息
	res, err := utils.SendMsg(c, producer, topic, tag, data)
	if err != nil {
		logger.SugarLogger.Error("发送消息失败， error:", err.Error())
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "操作失败",
			})
			return
		}
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
				c.JSON(http.StatusOK, CommentActionResponse{
					StatusCode: (int32)(code),
					StatusMsg:  info["status_msg"],
				})
				return
			}

			time.Sleep(time.Second / 15)
		}
	}
	// 20秒后，还是没有结果，则返回
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "请求超时",
		})
		return
	}
}
