package service

import (
	"simple_tiktok/dao/mysql"
	"simple_tiktok/models"

	"github.com/gin-gonic/gin"
)

// 请求体
type CommentActionRequire struct {
	Model 		models.CommentVideo
	ActionType  int
}

func PostCommentVideoAction(c *gin.Context, req *CommentActionRequire) error {
	if(req.ActionType == 1) {
		// 发表评论
		err := mysql.AddComment(req.Model)
		if err != nil {
			return err
		}
	} else {
		// 删除评论
		err := mysql.DelComment(req.Model.Identity)
		if err != nil {
			return err
		}
	}
	return nil
}
