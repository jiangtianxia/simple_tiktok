package controller

import (
	"encoding/json"
	"fmt"
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
	Id              uint64 `json:"id"`
	Name            string `json:"name"`
	Follow_count    uint64 `json:"follow_count"`
	Follower_count  uint64 `json:"follower_count"`
	Is_follow       bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackGroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}

type RespComment struct {
	Id          uint64   `json:"id"`
	User        RespUser `json:"user"`
	Content     string   `json:"content"`
	Create_date string   `json:"create_date"`
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
	comment_text := c.DefaultQuery("comment_text", "")
	comment_id := c.DefaultQuery("comment_id", "0")

	// 参数处理
	// 错误返回空用户
	failresp := CommentActionResponse{}
	failresp.Comment.Id = 0
	failresp.Comment.Content = comment_text
	failresp.Comment.Create_date = ""
	failresp.Comment.User.Id = 0
	failresp.Comment.User.Name = ""
	failresp.Comment.User.Follow_count = (uint64)(0)
	failresp.Comment.User.Follower_count = (uint64)(0)
	failresp.Comment.User.Is_follow = false

	// 校验参数
	if action_type != "1" && action_type != "2" {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "非法参数",
			Comment:    failresp.Comment,
		})
		return
	}

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

	if actiontype == 1 {
		identity, err := utils.GetID()
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "操作失败",
				Comment:    failresp.Comment,
			})
			return
		}
		commentidentity = int(identity)
	}

	// 处理时间格式
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	timeNow := month + "-" + day

	// 数据打包，将消息发送到消息队列
	pack := models.CommentVideo{
		Identity:      (uint64)(commentidentity),
		VideoIdentity: (uint64)(videoidentity),
		UserIdentity:  user.Identity,
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

				if code == -1 {
					c.JSON(http.StatusOK, CommentActionResponse{
						StatusCode: int32(code),
						StatusMsg:  info["status_msg"],
						Comment:    failresp.Comment,
					})
					return
				}

				//给返回题赋值
				authorid, _ := mysql.SearchAuthorIdByVideoId((uint64)(videoidentity))
				followcount, _ := service.GetFollowCount(c, fmt.Sprintf("%d", authorid))
				followercount, _ := service.GetFollowerCount(c, fmt.Sprintf("%d", authorid))
				isfollow, _ := service.IsFollow(c, fmt.Sprintf("%d", authorid), fmt.Sprintf("%d", user.Identity))
				// 获取点赞数量，作品数和喜欢数
				totalFavourited, workCount, FavouriteCount, _ := service.GetTotalFavouritedANDWorkCountANDFavoriteCount(authorid)

				resp := CommentActionResponse{}
				resp.Comment.Id = (uint64)(commentidentity)
				resp.Comment.Content = comment_text
				resp.Comment.Create_date = timeNow
				resp.Comment.User.Id = user.Identity
				resp.Comment.User.Name = user.Username
				resp.Comment.User.Follow_count = (uint64)(followcount)
				resp.Comment.User.Follower_count = (uint64)(followercount)
				resp.Comment.User.Is_follow = isfollow
				resp.Comment.User.Avatar = viper.GetString("defaultAvatarUrl")
				resp.Comment.User.BackGroundImage = viper.GetString("defaultBackGroudImage")
				resp.Comment.User.Signature = viper.GetString("defaultSignature")
				resp.Comment.User.TotalFavorited = totalFavourited
				resp.Comment.User.WorkCount = workCount
				resp.Comment.User.FavoriteCount = FavouriteCount

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
