package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"simple_tiktok/dao/mysql"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"
	"time"
)

// 发送消息接收参数结构体
type SendMessageReqStruct struct {
	FromUserId uint64
	ToUserId   string
	ActionType string
	Content    string
}

// 聊天记录返回结构体
type MessageList struct {
	Identity   uint64 `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

func MessageRecord(c *gin.Context, fromUserId uint64, toUserId string) ([]MessageList, error) {
	setKey := strconv.FormatUint(fromUserId, 10) + viper.GetString("redis.KeyUserMessageListPrefix") + toUserId

	if utils.RDB12.Exists(c, setKey).Val() == 0 {
		// 类型转化
		toIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		toIdUint64 := uint64(toIdInt64)

		// 查询数据库
		res, err := mysql.QueryMessageByToUserIdentity(fromUserId, toIdUint64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		messageList := make([]MessageList, len(res))
		for i := 0; i < len(res); i++ {
			messageList[i].Identity = res[i].Identity
			messageList[i].Content = res[i].Content
			messageList[i].CreateTime = res[i].CreateTime
		}

		// 新增缓存
		err = myRedis.RedisAddUserMessageSet(setKey, res)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		err = myRedis.RedisAddUserMessageHash(res)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}

		fmt.Println("数据库")
		return messageList, nil
	}
	// 使用缓存
	cathe := utils.RDB12.ZRange(c, setKey, 0, utils.RDB12.ZCard(c, setKey).Val()).Val()
	messageList := make([]MessageList, len(cathe))

	hashKey := viper.GetString("redis.KeyUserMessageHashPrefix")
	for i := 0; i < len(cathe); i++ {
		// 判断缓存哈希表中是否有记录
		if utils.RDB13.Exists(c, hashKey+cathe[i]).Val() == 0 {
			// 累心转换
			idInt64, err := strconv.ParseInt(cathe[i], 10, 64)
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, err
			}
			idUint64 := uint64(idInt64)

			// 查询数据库
			res, err := mysql.QueryMessageByIdentity(idUint64)
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, err
			}

			// 新增缓存
			myRedis.RedisAddUserMessageHash(res)
			fmt.Println("补漏")
		}
		message := utils.RDB13.HGetAll(c, hashKey+cathe[i]).Val()
		idInt64, err := strconv.ParseInt(message["id"], 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return nil, err
		}
		idUint64 := uint64(idInt64)
		messageList[i].Identity = idUint64
		messageList[i].Content = message["content"]
		messageList[i].CreateTime = message["create_time"]
	}

	fmt.Println("缓存")
	return messageList, nil
}

func SendMessage(msgid string, data []byte) {
	SendMessageReqStruct := SendMessageReqStruct{}
	json.Unmarshal(data, &SendMessageReqStruct)

	fromUserId := SendMessageReqStruct.FromUserId
	toUserId := SendMessageReqStruct.ToUserId
	actionType := SendMessageReqStruct.ActionType
	content := SendMessageReqStruct.Content

	if actionType == "1" {
		// 雪花算法生成identity
		identity, err := utils.GetID()
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "操作失败")
			return
		}

		// 类型转化
		toIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "操作失败")
			return
		}
		toIdUint64 := uint64(toIdInt64)

		userMessage := models.UserMessage{
			Identity:         identity,
			ToUserIdentity:   toIdUint64,
			FromUserIdentity: fromUserId,
			Content:          content,
			CreateTime:       time.Now().Format("2006-01-02 15:04:05"),
		}

		// 存入数据库
		err = mysql.CreateUserMessage(userMessage)
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "操作失败")
			return
		}

		// 发送延迟消息，删除缓存
		RetryTopic := viper.GetString("rocketmq.RetryTopic")
		DeleteMssageRedisTag := viper.GetString("rocketmq.DeleteMessageRedisTag")
		utils.SendDelayMsg(RetryTopic, DeleteMssageRedisTag, data)
		SaveRedisResp(msgid, 0, "操作成功")
		return
	}
	SaveRedisResp(msgid, -1, "操作失败")
}

// 将结果存入redis缓存
func SaveRedisResp(msgid string, code int, msg string) {
	info := map[string]interface{}{
		"status_code": code,
		"status_msg":  msg,
	}

	var ctx = context.Background()
	pipeline := utils.RDB0.Pipeline()
	pipeline.HSet(ctx, msgid, info)
	pipeline.Expire(ctx, msgid, time.Second*70)
	pipeline.Exec(ctx)
}
