package rocket

import (
	"fmt"
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

var FollowChan chan ChanMsg = make(chan ChanMsg, 100)

func ReceiveChan() {
	for {
		select {
		case data := <-FollowChan:
			// 关注操作
			service.FollowService(data.Msgid, data.Data)
		}
	}
}
