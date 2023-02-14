package service

import (
	"encoding/json"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/models"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: service层，接受controller层参数-调用dao层函数
 **/

// 请求体
type CommentActionRequire struct {
	Model 		models.CommentVideo
	ActionType  int
}

func PostCommentVideoAction(msgid string, data []byte) {
	req := &CommentActionRequire{}
	json.Unmarshal(data, req)
	if(req.ActionType == 1) {
		// 发表评论
		err := mysql.AddComment(req.Model)
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			return
		}
	} else {
		// 删除评论
		err := mysql.DelComment(req.Model.Identity)
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			return
		}
	}
	return
}
