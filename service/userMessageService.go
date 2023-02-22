package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"simple_tiktok/dao/mysql"
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

/**
 * @Author
 * @Description 发送消息
 * @Date 11:00 2023/2/14
 **/
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
			SaveRedisResp(msgid, -1, "发送失败")
			return
		}

		// 类型转化
		toIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "发送失败")
			return
		}
		toIdUint64 := uint64(toIdInt64)

		userMessage := models.UserMessage{
			Identity:         identity,
			ToUserIdentity:   toIdUint64,
			FromUserIdentity: fromUserId,
			Content:          content,
			CreateTime:       time.Now().Unix(),
		}

		// 先存缓存
		setKey := strconv.FormatUint(fromUserId, 10) + viper.GetString("redis.KeyUserMessageListPrefix") + toUserId
		pipeline := utils.RDB12.Pipeline()
		var ctx = context.Background()
		pipeline.RPush(ctx, setKey, userMessage.Identity)
		hashKey := viper.GetString("redis.KeyUserMessageHashPrefix") + fmt.Sprintf("%d", userMessage.Identity)
		value := map[string]interface{}{
			"id":          userMessage.Identity,
			"content":     userMessage.Content,
			"create_time": userMessage.CreateTime,
		}
		pipeline.HSet(ctx, hashKey, value)
		pipeline.Expire(ctx, hashKey, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
		_, err = pipeline.Exec(ctx)
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "发送失败")
			return
		}

		// 存入数据库
		err = mysql.CreateUserMessage(userMessage)
		if err != nil {
			logger.SugarLogger.Error(err)
			SaveRedisResp(msgid, -1, "发送失败")
			return
		}

		SaveRedisResp(msgid, 0, "发送成功")
		return
	}
	SaveRedisResp(msgid, -1, "发送失败")
}

/**
 * @Author
 * @Description 聊天记录
 * @Date 11:00 2023/2/14
 **/
func MessageRecord(c *gin.Context, fromUserId uint64, toUserId string) ([]Message, error) {
	setKey := toUserId + viper.GetString("redis.KeyUserMessageListPrefix") + strconv.FormatUint(fromUserId, 10)

	len, err := utils.RDB12.LLen(c, setKey).Result()
	if err != nil {
		logger.SugarLogger.Error(err)
		return []Message{}, err
	}

	messageList := make([]Message, 0)
	for i := 0; i < int(len); i++ {
		// 出队列
		messageId := utils.RDB12.LPop(c, setKey).Val()

		// 查询缓存
		hashKey := viper.GetString("redis.KeyUserMessageHashPrefix") + messageId
		// 判断缓存哈希表中是否有记录
		if utils.RDB12.Exists(c, hashKey).Val() == 0 {
			// 类型转换
			idInt64, _ := strconv.ParseInt(messageId, 10, 64)

			// 查询数据库
			res, err := mysql.QueryMessageByIdentity(uint64(idInt64))
			if err != nil {
				logger.SugarLogger.Error(err)
				return []Message{}, err
			}

			Message := Message{
				Identity:   res.Identity,
				Content:    res.Content,
				CreateTime: res.CreateTime,
				FromUserId: int64(res.FromUserIdentity),
				ToUserId:   int64(res.ToUserIdentity),
			}
			messageList = append(messageList, Message)

			// 新增缓存
			pipeline := utils.RDB12.Pipeline()
			value := map[string]interface{}{
				"id":          res.Identity,
				"content":     res.FromUserIdentity,
				"create_time": res.CreateTime,
			}
			pipeline.HSet(c, hashKey, value)
			pipeline.Expire(c, hashKey, time.Hour*time.Duration(viper.GetInt("redis.RedisExpireTime")))
			pipeline.Exec(c)
		} else {
			// 使用缓存
			message := utils.RDB12.HGetAll(c, hashKey).Val()
			idInt64, _ := strconv.ParseInt(message["id"], 10, 64)
			to_user_id, _ := strconv.ParseInt(toUserId, 10, 64)
			create_time, _ := strconv.ParseInt(message["create_time"], 10, 64)
			Message := Message{
				Identity:   uint64(idInt64),
				Content:    message["content"],
				CreateTime: create_time,
				FromUserId: to_user_id,
				ToUserId:   int64(fromUserId),
			}
			messageList = append(messageList, Message)
			utils.RDB12.Del(c, hashKey)
		}
	}
	return messageList, nil
}
