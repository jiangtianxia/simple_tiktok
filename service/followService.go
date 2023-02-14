package service

import (
	"encoding/json"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"

	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 关注操作
 * @Date 16:00 2023/2/11
 **/
/**
 * 修改方法
 * 1、先删除缓存，后更新数据库
 * 2、先更新数据库，后删除缓存
 *
 * 保证数据库数据和缓存数据一致性方法
 * 先删除缓存，后更新数据库：延时双删。
 * 先更新数据库，后删除缓存：重试机制或binlog日志。
 *
 * 采用：先更新数据库，后删缓存+重试机制。
 *
 * 为什么用重试机制？
 * 答：删除缓存时可能会删除失败从而导致数据不一致，因此可以用重试机制，重试删除缓存，如果次数过多需要人工介入。
 *
 * 如何实现重试机制？
 * 答：使用rocketmq来实现重试机制。
 **/
func FollowService(msgid string, data []byte) {
	followInfo := &FollowReqStruct{}
	json.Unmarshal(data, followInfo)
	// fmt.Println("userid：", followInfo.UserId, "toUserId：", followInfo.ToUserId, "type：", followInfo.ActionType)

	// 1、更新数据库
	userid, _ := strconv.Atoi(followInfo.UserId)
	status := 0
	if followInfo.ActionType == 1 {
		status = 1
	}
	followId, _ := strconv.Atoi(followInfo.ToUserId)
	err := UpdateUserFollow(uint64(userid), uint64(followId), status)
	if err != nil {
		logger.SugarLogger.Error("UpdateUserFollow Error：", err.Error())
		SaveRedisResp(msgid, -1, "操作失败")
		return
	}

	// 2、发送延迟消息，删除缓存
	RetryTopic := viper.GetString("rocketmq.RetryTopic")
	DeleteFollowRedisTag := viper.GetString("rocketmq.DeleteFollowRedisTag")
	err = utils.SendDelayMsg(RetryTopic, DeleteFollowRedisTag, data)
	if err != nil {
		logger.SugarLogger.Error("SendDelayMsg Error：", err.Error())
		// SaveRedisResp(msgid, -1, "操作失败")
		return
	}
	SaveRedisResp(msgid, 0, "操作成功")
}

// 更新关注表信息
func UpdateUserFollow(identity uint64, followIdentity uint64, status int) error {
	// 查询关注表中是否存在数据
	cnt, err := mysql.FindUserFollowByIdentityCount(identity, followIdentity)
	if err != nil {
		return err
	}

	info := models.UserFollow{
		UserIdentity:     identity,
		FollowerIdentity: followIdentity,
		Status:           status,
	}
	if cnt > 0 {
		// 存在，则修改
		return mysql.UpdateUserFollow(info)
	}
	// 不存在，则创建
	return mysql.CreateUserFollow(info)
}
