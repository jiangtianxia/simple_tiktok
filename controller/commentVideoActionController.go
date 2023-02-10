package controller

import (
	"net/http"
	//"simple_tiktok/middlewares"
	"simple_tiktok/models"
	"simple_tiktok/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 返回体
type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    models.CommentVideo `json:"comment"`
}

// 发表，删除评论 /comment/action
func CommentAction(c *gin.Context) {
	//参数获取
	//token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	comment_text := c.Query("comment_text")
	comment_id := c.Query("comment_id")

	//参数处理
	//验证用户token
	//user, err := middlewares.AuthUserCheck(token)


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
	timeNow := time.Now()

	// 数据打包
	pack := models.CommentVideo{
		VideoIdentity: (uint64)(videoidentity),
		UserIdentity:  (uint64)(commentidentity),
		Text:          comment_text,
		CommentTime:   timeNow.String(),
	}
	req := service.CommentActionRequire{
		Model:      pack,
		ActionType: actiontype,
	}

	//调用service层函数
	err = service.PostCommentVideoAction(c, &req)

	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg: err.Error(),
		})
		return
	}
	if actiontype == 1 {
		//发表评论
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg: "发表评论成功",
		})
		return
	} else {
		//删除评论
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg: "删除评论成功",
		})
		return
	}
}
