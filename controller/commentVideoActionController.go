package controller

import (
	"encoding/json"
	"net/http"

	"simple_tiktok/dao/mysql"
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
type RespUser struct {
	id             uint64 `json:"id"`
	name           string `json:"name"`
	follow_count   uint64 `json:"follow_count"`
	follower_count uint64 `json:"follower_count"`
	is_follow      bool   `json:"is_follow"`
}

type RespComment struct {
	id          uint64   `json:"id"`
	user        RespUser `json:"user"`
	content     string   `json:"content"`
	create_date string   `json:"create_date"`
}

type CommentActionResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	Comment    RespComment `json:"comment"`
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
	// 错误返回空用户
	failresp := CommentActionResponse{}
	failresp.Comment.id = (uint64)(0)
	failresp.Comment.content = comment_text
	failresp.Comment.create_date = ""
	failresp.Comment.user.id = 0
	failresp.Comment.user.name = ""
	failresp.Comment.user.follow_count = (uint64)(0)
	failresp.Comment.user.follower_count = (uint64)(0)
	failresp.Comment.user.is_follow = false
	// 验证用户token
	user, err := middlewares.AuthUserCheck(token)
	if user == nil || user.Identity == 0 || user.Issuer != "simple_tiktok" || user.Username == "" || err != nil {
		logger.SugarLogger.Error("Unauthorized User")
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户错误",
			Comment:    failresp.Comment,
		})
		return
	}

	videoidentity, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "评论id格式错误",
			Comment:    failresp.Comment,
		})
		return
	}
	commentidentity, err := strconv.Atoi(comment_id)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频id格式错误",
			Comment:    failresp.Comment,
		})
		return
	}
	actiontype, err := strconv.Atoi(action_type)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "操作类型错误",
			Comment:    failresp.Comment,
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
				Comment:    failresp.Comment,
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
				//给返回题赋值
				followcount, _ := service.GetFollowCount(c, (string)(user.Identity))
				followercount, _ := service.GetFollowerCount(c, (string)(user.Identity))
				authorid, _ := mysql.SearchAuthorIdByVideoId((uint64)(commentidentity))
				isfollow, _ := service.IsFollow(c, (string)(user.Identity), (string)(authorid))

				resp := CommentActionResponse{}
				resp.Comment.id = (uint64)(commentidentity)
				resp.Comment.content = comment_text
				resp.Comment.create_date = timeNow
				resp.Comment.user.id = user.Identity
				resp.Comment.user.name = user.Username
				resp.Comment.user.follow_count = (uint64)(followcount)
				resp.Comment.user.follower_count = (uint64)(followercount)
				resp.Comment.user.is_follow = isfollow

				c.JSON(http.StatusOK, CommentActionResponse{
					StatusCode: (int32)(code),
					StatusMsg:  info["status_msg"],
					Comment:    resp.Comment,
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
			Comment:    failresp.Comment,
		})
		return
	}
}

/*
{
    "status_code": 0,
    "status_msg": "string",
    "comment": {
        "id": 0,
        "user": {
            "id": 0,
            "name": "string",
            "follow_count": 0,
            "follower_count": 0,
            "is_follow": true
        },
        "content": "string",
        "create_date": "string"
    }
}
*/
