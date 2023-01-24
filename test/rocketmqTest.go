package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	// // 创建主题
	// CreateTopic("SimpleTopic")
	// fmt.Println("创建主题成功")

	// 发送消息
	for i := 0; i < 100; i++ {
		SendMsg()
	}

	// // 接收消息
	// for {
	// 	ReceiveMsg()
	// }
}

func CreateTopic(topicName string) {
	endPoint := []string{"192.168.1.69:9876"}
	// 创建主题
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(endPoint)))
	if err != nil {
		fmt.Printf("connection error: %s\n", err.Error())
	}
	err = testAdmin.CreateTopic(context.Background(),
		admin.WithTopicCreate(topicName),
		admin.WithBrokerAddrCreate("192.168.1.69:10909"),
	)
	if err != nil {
		fmt.Printf("createTopic error: %s\n", err.Error())
	}
}

// 发送消息
func SendMsg() {
	// 发送消息
	endPoint := []string{"127.0.0.1:9876"}

	// 消息消费失败重试两次
	newProducer, err := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(2),
		producer.WithQueueSelector(producer.NewRandomQueueSelector()),
	)

	defer func(newProducer rocketmq.Producer) {
		err := newProducer.Shutdown()
		if err != nil {
			panic("关闭producer失败")
		}
	}(newProducer)
	if err != nil {
		panic("生成producer失败")
	}

	if err = newProducer.Start(); err != nil {
		panic("启动producer失败")
	}

	msg := primitive.NewMessage("SimpleTopic", []byte("大家好"))
	msg.WithTag("test")
	res, err := newProducer.SendSync(context.Background(), msg)

	if err != nil {
		panic("消息发送失败" + err.Error())
	}
	fmt.Println("***************")
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s: 消息: %s发送成功 \n", nowStr, res.String())
}

// 接收消息
func ReceiveMsg() {
	// 发送消息
	endPoint := []string{"127.0.0.1:9876"}

	// 消息消费失败重试两次
	newPushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(endPoint),
		consumer.WithGroupName("Test"),
		consumer.WithPullBatchSize(2),
		consumer.WithConsumeMessageBatchMaxSize(2),
	)
	if err != nil {
		fmt.Println("创建消费者失败")
	}
	defer func(newPushConsumer rocketmq.PushConsumer) {
		err := newPushConsumer.Shutdown()
		if err != nil {
			panic("关闭consumer失败")
		}
	}(newPushConsumer)

	err = newPushConsumer.Subscribe("SimpleTopic", consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			fmt.Println("1")
			for _, msg := range msgs {
				nowStr := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("%s 读取到一条消息,消息内容: %s Tags: %s msgId: %s \n", nowStr, string(msg.Body), msg.GetTags(), msg.MsgId)
				fmt.Println("睡眠")

				time.Sleep(time.Second * 10)
			}
			fmt.Println("2")
			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		fmt.Println("读取消息失败")
	}
	if err = newPushConsumer.Start(); err != nil {
		panic("启动consumer失败")
	}

	// 不能让主goroutine退出
	time.Sleep(time.Second * 60)
}
