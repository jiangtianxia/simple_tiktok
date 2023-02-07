package service

import (
	"errors"
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

// 聊天记录返回结构体
type MessageList struct {
	Identity   uint64 `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

func SendMessage(c *gin.Context, fromUserId uint64, toUserId string, actionType string, content string) error {
	if actionType == "1" {
		// 雪花算法生成identity
		identity, err := utils.GetID()
		if err != nil {
			logger.SugarLogger.Error(err)
			return err
		}

		// 类型转化
		toIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
		if err != nil {
			logger.SugarLogger.Error(err)
			return err
		}
		toIdUint64 := uint64(toIdInt64)

		userMessage := models.UserMessage{
			Identity:         identity,
			ToUserIdentity:   toIdUint64,
			FromUserIdentity: fromUserId,
			Content:          content,
			CreateTime:       time.Now().Format("2006-01-02 15:04:05"),
		}

		// 删除缓存
		setKey := strconv.FormatUint(fromUserId, 10) + viper.GetString("redis.KeyUserMessageListPrefix") + toUserId
		utils.RDB12.Del(c, setKey)

		// 存入数据库
		err = mysql.CreateUserMessage(userMessage)
		if err != nil {
			logger.SugarLogger.Error(err)
			return err
		}

		// 延时双删
		time.Sleep(1 * time.Second)
		utils.RDB12.Del(c, setKey)

		return nil
	}
	return errors.New("发送消息失败")
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
