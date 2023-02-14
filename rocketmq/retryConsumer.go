package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 重试机制
 * @Date 20:00 2023/2/11
 **/
func CreateDelayConsumer(groupName string, topic string, tags string) {
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
		ReceiveDelayMsg(newPushConsumer, topic, tags)
	}
}

func ReceiveDelayMsg(newPushConsumer rocketmq.PushConsumer, topic string, tags string) {
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

				// 判断当前消息重试次数
				if msg.ReconsumeTimes > 3 {
					// 如果重试次数大于3，则返回成功，此时需要人工介入，此处不设人工处理
					return consumer.ConsumeSuccess, nil
				}

				/*
				* 选择器
				 */
				switch msg.GetTags() {
				case "DeleteFollowRedis":
					// 删除缓存
					// fmt.Println("接收到消息，重试次数：", msg.ReconsumeTimes)
					followInfo := &FollowReqStruct{}
					json.Unmarshal(msg.Body, followInfo)

					err := redis.DeleteFollowList(followInfo.UserId)
					if err != nil {
						logger.SugarLogger.Error("DeleteFollowList Error：", err.Error())
						return consumer.ConsumeRetryLater, nil
					}

					err = redis.DeleteFollowerSortSet(followInfo.ToUserId)
					if err != nil {
						logger.SugarLogger.Error("DeleteFollowerSortSet Error：", err.Error())
						return consumer.ConsumeRetryLater, nil
					}
					// fmt.Println("消息执行成功")
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

// 接收参数结构体
type FollowReqStruct struct {
	UserId     string
	ToUserId   string
	ActionType int
}
