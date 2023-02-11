package rocket

import (
	"encoding/json"
	"fmt"
	"simple_tiktok/models"
	"simple_tiktok/service"

	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description rocketmq初始化
 * @Date 21:00 2023/1/22
 **/

func InitRocketmq() {
	// 打开通道协程，一个通道接收一个接口的数据
	go ReceiveChan()

	// 创建消费者组
	redisTopic := viper.GetString("rocketmq.redisTopic")
	redisGroupName := viper.GetString("rocketmq.redisGroupName")
	redisTags := viper.GetString("rocketmq.redisTags")

	go CreateConsumer(redisGroupName, redisTopic, redisTags)

	// 关注操作的消费者组
	followTopic := viper.GetString("rocketmq.ServerTopic")
	followTag := viper.GetString("rocketmq.serverFollowTag")
	followGroupName := viper.GetString("rocketmq.followGroupName")
	go CreateConsumer(followGroupName, followTopic, followTag)

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

var publishChan chan ChanMsg = make(chan ChanMsg, 100)
var LoginChan chan ChanMsg = make(chan ChanMsg, 100)
var FollowChan chan ChanMsg = make(chan ChanMsg, 100)

func ReceiveChan() {
	for {
		select {
		case data := <-publishChan:
			// 用户上传视频时，发送videobasic到消息队列，将信息缓存到redis
			videoinfo := &models.VideoBasic{}
			json.Unmarshal(data.Data, videoinfo)
			// fmt.Println(videoinfo)
			PublishAction(*videoinfo)
		case data := <-FollowChan:
			// 关注操作
			service.FollowService(data.Msgid, data.Data)
		}
	}
}
