package controller

import (
	"net/http"
	"simple_tiktok/models"
	"simple_tiktok/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 请求体
type CommentActionRequire struct {
	Model models.CommentVideo
}

type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    service.CommentActionService `json:"comment"`
}

// 发表，删除评论 comment/action/
func CommentAction(c *gin.Context) {
	//参数获取
	//token := c.Query("token")
	video_id:= c.Query("video_id")
	action_type:= c.Query("action_type")
	comment_text:= c.Query("comment_text")
	comment_id:= c.Query("comment_id")
	
	//参数处理
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
	timeNow := time.Now()


	//调用service层函数
	if action_type == "1" { 
		// 发表评论
		
		// 数据打包
		var newComment models.CommentVideo
		newComment.VideoIdentity = (uint64)(videoidentity)
		newComment.UserIdentity = (uint64)(commentidentity)
		newComment.Text = comment_text
		newComment.CreatedAt = timeNow
	
		//发表评论
		commentInfo, err := service.AddComment(newComment)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "发表评论失败",
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "发表评论成功",
			Comment:    commentInfo,
		})
		return

	} else { 
		//删除评论
		//service ->缺参数
		//commentInfo, err := service.DelComment() 
		if err != nil { 
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "删除评论失败",
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "删除评论成功",
		})
		return
	}
}
