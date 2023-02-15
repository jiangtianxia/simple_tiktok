package service

import (
	"encoding/json"

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
	setKey := strconv.FormatUint(fromUserId, 10) + viper.GetString("redis.KeyUserMessageListPrefix") + toUserId

	if utils.RDB12.Exists(c, setKey).Val() == 0 {
		// 类型转化
		toIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return []Message{}, err
		}
		toIdUint64 := uint64(toIdInt64)

		// 查询数据库
		res, err := mysql.QueryMessageByToUserIdentity(fromUserId, toIdUint64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return []Message{}, err
		}
		messageList := make([]Message, len(res))
		for i := 0; i < len(res); i++ {
			messageList[i].Identity = res[i].Identity
			messageList[i].Content = res[i].Content
			messageList[i].CreateTime = res[i].CreateTime
		}

		// 新增缓存
		err = myRedis.RedisAddUserMessageSet(c, setKey, res)
		if err != nil {
			logger.SugarLogger.Error(err)
			return []Message{}, err
		}
		err = myRedis.RedisAddUserMessageHash(c, res)
		if err != nil {
			logger.SugarLogger.Error(err)
			return []Message{}, err
		}

		return messageList, nil
	}
	// 使用缓存
	cathe := utils.RDB12.ZRange(c, setKey, 0, utils.RDB12.ZCard(c, setKey).Val()).Val()
	messageList := make([]Message, len(cathe))

	hashKey := viper.GetString("redis.KeyUserMessageHashPrefix")
	for i := 0; i < len(cathe); i++ {
		// 判断缓存哈希表中是否有记录
		if utils.RDB13.Exists(c, hashKey+cathe[i]).Val() == 0 {
			// 累心转换
			idInt64, err := strconv.ParseInt(cathe[i], 10, 64)
			if err != nil {
				logger.SugarLogger.Error(err)
				return []Message{}, err
			}
			idUint64 := uint64(idInt64)

			// 查询数据库
			res, err := mysql.QueryMessageByIdentity(idUint64)
			if err != nil {
				logger.SugarLogger.Error(err)
				return []Message{}, err
			}

			// 新增缓存
			myRedis.RedisAddUserMessageHash(c, res)
			// fmt.Println("补漏")
		}
		message := utils.RDB13.HGetAll(c, hashKey+cathe[i]).Val()
		idInt64, err := strconv.ParseInt(message["id"], 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return []Message{}, err
		}
		idUint64 := uint64(idInt64)
		messageList[i].Identity = idUint64
		messageList[i].Content = message["content"]
		creareTime, _ := strconv.Atoi(message["create_time"])
		messageList[i].CreateTime = int64(creareTime)
	}

	// fmt.Println("缓存")
	return messageList, nil
}
