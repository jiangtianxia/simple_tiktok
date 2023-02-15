package service

import (
	"encoding/json"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"

	"github.com/spf13/viper"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: service层，接受controller层参数-调用dao层函数
 **/

// 请求体
type CommentActionRequire struct {
	Model      models.CommentVideo
	ActionType int
}

func PostCommentVideoAction(msgid string, data []byte) {
	req := &CommentActionRequire{}
	json.Unmarshal(data, req)
	// 更新数据库
	if req.ActionType == 1 {
		// identity, err := utils.GetID()
		// if err != nil {
		// 	logger.SugarLogger.Error(err.Error())
		// 	SaveRedisResp(msgid, -1, "操作失败")
		// }
		// // 发表评论
		// req.Model.Identity = identity
		err := mysql.AddComment(req.Model)
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			SaveRedisResp(msgid, -1, "操作失败")
			return
		}
	}

	if req.ActionType == 2 {
		// 删除评论
		err := mysql.DelComment(req.Model.Identity)
		if err != nil {
			logger.SugarLogger.Error(err.Error())
			SaveRedisResp(msgid, -1, "操作失败")
			return
		}
	}

	//发送延迟消息，删除缓存
	RetryTopic := viper.GetString("rocketmq.RetryTopic")
	DeleteFollowRedisTag := viper.GetString("rocketmq.DeleteCommentRedisTag")
	err := utils.SendDelayMsg(RetryTopic, DeleteFollowRedisTag, data)
	if err != nil {
		logger.SugarLogger.Error("SendDelayMsg Error：", err.Error())
	}
	//将结果存入redis缓存
	SaveRedisResp(msgid, 0, "操作成功")
}
