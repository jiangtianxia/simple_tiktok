package rocket

import (
	"context"
	"fmt"
	"simple_tiktok/logger"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 消费者组接收消息
 * @Date 15:00 2023/1/23
 **/
/**
 * 此次共有两个消息费者组。
 * 1）服务的消息者，获取服务消息队列的数据，同时根据tag进行筛选。
 * 2）缓存的消息者，获取缓存消息队列的数据，同时根据tag进行筛选。
 * 注：为了防止数据丢失，一个消费者组接收一个topic的消息
 **/
func CreateConsumer(groupName string, topic string, tags string) {
	// 服务地址
	endPoint := []string{viper.GetString("rocketmq.addr")}

	// 每次拉取16条消息，每次处理16条消息，每隔0.5秒钟拉取一次
	newPushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(endPoint),
		consumer.WithGroupName(groupName),
		consumer.WithPullBatchSize(16),
		consumer.WithConsumeMessageBatchMaxSize(16),
		consumer.WithPullInterval(time.Second/2),
	)

	if err != nil {
		fmt.Println("创建消费者失败，error：", err.Error())
		logger.SugarLogger.Error("创建消费者失败，error：", err.Error())
	}
	defer func(newPushConsumer rocketmq.PushConsumer) {
		err := newPushConsumer.Shutdown()
		if err != nil {
			fmt.Println("关闭消费者失败，error：", err.Error())
			logger.SugarLogger.Error("关闭消费者失败，error：", err.Error())
		}
	}(newPushConsumer)

	// 接收消息
	for {
		ReceiveMsg(newPushConsumer, topic, tags)
	}
}

func ReceiveMsg(newPushConsumer rocketmq.PushConsumer, topic string, tags string) {
	// 过滤器，只接收主题为topic，标签为tag的数据
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tags,
	}
	err := newPushConsumer.Subscribe(topic, selector,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				nowStr := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("%s 读取到一条消息,消息内容: %s Tags: %s msgId: %s \n", nowStr, string(msg.Body), msg.GetTags(), msg.MsgId)

				/*
				* 选择器，根据tag判断要将数据发送至哪条通道
				 */
				switch msg.GetTags() {
				case "publishAction":
					// 把msgid，也发到到通道
					publishChan <- ChanMsg{
						Msgid: msg.MsgId,
						Data:  msg.Body,
					}
				case "loginredis":
					LoginChan <- ChanMsg{
						Msgid: msg.MsgId,
						Data:  msg.Body,
					}
				case "userInfo":
					userInfoChan <- ChanMsg{
						Msgid: msg.MsgId,
						Data:  msg.Body,
					}
				}
			}
			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		fmt.Println("读取消息失败，error：", err.Error())
		logger.SugarLogger.Error("读取消息失败，error：", err.Error())
	}
	if err = newPushConsumer.Start(); err != nil {
		fmt.Println("启动consumer失败，error：", err.Error())
		logger.SugarLogger.Error("启动consumer失败，error：", err.Error())
		return
	}

	// 不能让主goroutine退出
	time.Sleep(time.Hour * 24)
}
