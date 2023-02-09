package service

import (
	"simple_tiktok/controller"

	"github.com/gin-gonic/gin"
)

// CommentActionService 接口定义
// 发表评论-使用的结构体-service层引用dao层↑的Comment
func PostCommentVideoAction(c *gin.Context, req *controller.CommentActionRequire) (*controller.CommentActionResponse, error) {
	// 根据videoId获取视频评论数量
	//req  ->  model -> id
	//QueryVideoList(id uint64)  //mysql.QueryVideoList
	// 发表评论，传进来评论的基本信息，返回保存是否成功的状态描述
	//AddComment(comment models.CommentVideo) (models.CommentVideo, error)
	// 删除评论，传入评论id即可，返回错误状态信息
	//DelComment(commentId int64) error
}
