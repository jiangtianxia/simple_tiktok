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

	// 赞操作的消费者组
	FavouritTopic := viper.GetString("rocketmq.ServerTopic")
	FavouritTag := viper.GetString("rocketmq.serverFavouriteTag")
	FavouriteGroupName := viper.GetString("rocketmq.FavouriteGroupName")
	go CreateConsumer(FavouriteGroupName, FavouritTopic, FavouritTag)

	// 发表评论的消费者组
	CommentTopic := viper.GetString("rocketmq.ServerTopic")
	CommentTag := viper.GetString("rocketmq.serverSendCommentTag")
	CommentGroupName := viper.GetString("rocketmq.SendCommentGroupName")
	go CreateConsumer(CommentGroupName, CommentTopic, CommentTag)

	// 关注操作的消费者组
	followTopic := viper.GetString("rocketmq.ServerTopic")
	followTag := viper.GetString("rocketmq.serverFollowTag")
	followGroupName := viper.GetString("rocketmq.followGroupName")
	go CreateConsumer(followGroupName, followTopic, followTag)

	// 发送消息的消费者组
	messageTopic := viper.GetString("rocketmq.serverTopic")
	sendMessageGroupName := viper.GetString("rocketmq.sendMessageGroupName")
	sendMessageTags := viper.GetString("rocketmq.serverSendMessageTags")
	go CreateConsumer(sendMessageGroupName, messageTopic, sendMessageTags)

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
var sendMessageChan chan ChanMsg = make(chan ChanMsg, 100)
var commentActionChan chan ChanMsg = make(chan ChanMsg, 100)
var favouriteChan chan ChanMsg = make(chan ChanMsg, 100)

func ReceiveChan() {
	for {
		select {
		case data := <-favouriteChan:
			// 赞操作
			go service.DealFavourite(data.Msgid, data.Data)
		case data := <-commentActionChan:
			// 发表评论
			go service.PostCommentVideoAction(data.Msgid, data.Data)
		case data := <-FollowChan:
			// 关注操作
			go service.FollowService(data.Msgid, data.Data)
		case data := <-sendMessageChan:
			// 发送消息
			go service.SendMessage(data.Msgid, data.Data)
		}
	}
}
