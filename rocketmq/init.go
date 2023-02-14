package rocket

import (
	"fmt"
<<<<<<< HEAD
	"simple_tiktok/service"

	"github.com/spf13/viper"
=======
	"github.com/spf13/viper"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/service"
>>>>>>> zxy
)

/**
 * @Author jiang
 * @Description rocketmq初始化
 * @Date 21:00 2023/1/22
 **/

func InitRocketmq() {
	// 打开通道协程，一个通道接收一个接口的数据
	go ReceiveChan()

	// 关注操作的消费者组
	followTopic := viper.GetString("rocketmq.ServerTopic")
	followTag := viper.GetString("rocketmq.serverFollowTag")
	followGroupName := viper.GetString("rocketmq.followGroupName")
	go CreateConsumer(followGroupName, followTopic, followTag)

	// 重试机制
	RetryTopic := viper.GetString("rocketmq.RetryTopic")
	RetryTags := viper.GetString("rocketmq.RetryTags")
	RetryGroupName := viper.GetString("rocketmq.RetryGroupName")
	go CreateDelayConsumer(RetryGroupName, RetryTopic, RetryTags)

<<<<<<< HEAD
=======
	go CreateConsumer(redisGroupName, redisTopic, redisTags)

	// 发送消息的消费者组
	messageTopic := viper.GetString("rocketmq.serverTopic")
	sendMessageGroupName := viper.GetString("rocketmq.sendMessageGroupName")
	sendMessageTags := viper.GetString("rocketmq.serverSendMessageTags")

	go CreateConsumer(sendMessageGroupName, messageTopic, sendMessageTags)

	// 重试机制
	RetryTopic := viper.GetString("rocketmq.RetryTopic")
	// DeleteFollowRedisTag := viper.GetString("rocketmq.DeleteFollowRedisTag")
	RetryTags := viper.GetString("rocketmq.RetryTags")
	RetryGroupName := viper.GetString("rocketmq.RetryGroupName")
	go CreateDelayConsumer(RetryGroupName, RetryTopic, RetryTags)
>>>>>>> zxy
	fmt.Println("rocketmq inited ...... ")
}

/**
 * @Author jiang
 * @Description 通道协程
 * @Date 13:30 2023/1/23
 **/
// 通道，用于接收数据，一个通道接收一个接口的数据
type ChanMsg struct {
	Msgid string
	Data  []byte
}

<<<<<<< HEAD
var FollowChan chan ChanMsg = make(chan ChanMsg, 100)
=======
var publishChan chan ChanMsg = make(chan ChanMsg, 100)
var LoginChan chan ChanMsg = make(chan ChanMsg, 100)
var userInfoChan chan ChanMsg = make(chan ChanMsg, 100)
var sendMessageChan chan ChanMsg = make(chan ChanMsg, 100)
>>>>>>> zxy

func ReceiveChan() {
	for {
		select {
<<<<<<< HEAD
		case data := <-FollowChan:
			// 关注操作
			service.FollowService(data.Msgid, data.Data)
=======
		case data := <-publishChan:
			// 用户上传视频时，发送videobasic到消息队列，将信息缓存到redis
			videoinfo := &models.VideoBasic{}
			json.Unmarshal(data.Data, videoinfo)
			// fmt.Println(videoinfo)
			PublishAction(*videoinfo)
		case data := <-userInfoChan:
			userInfo := &models.UserBasic{}
			err := json.Unmarshal(data.Data, userInfo)
			if err != nil {
				logger.SugarLogger.Error(err)
			}
			UserInfoAction(*userInfo)
		case data := <-sendMessageChan:
			go service.SendMessage(data.Msgid, data.Data)
>>>>>>> zxy
		}
	}
}
