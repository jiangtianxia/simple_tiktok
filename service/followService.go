package service

import (
	"encoding/json"
)

// 接收参数结构体
type FollowReqStruct struct {
	UserId     string
	ToUserId   string
	ActionType int
}

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
 * 答：使用rocketmq来实现重试机制，发送延迟消息。
 **/
func FollowService(msgid string, data []byte) {
	followInfo := &FollowReqStruct{}
	json.Unmarshal(data, followInfo)
	// fmt.Println("userid：", followInfo.UserId, "toUserId：", followInfo.ToUserId, "type：", followInfo.ActionType)

	// 更新数据库

}
