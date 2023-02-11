package utils

import (
	"context"
	"fmt"
	"simple_tiktok/logger"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 生产者发送消息
 * @Date 13:30 2023/1/23
 **/
/**
 * 此次共有两个消息队列。
 * 1）服务的消息队列，接收访问接口的数据，不同服务使用不同tag标记。
 * 2）缓存的消息队列，接收将要加入缓存的数据，不同的缓存策略使用不同tag标记。
 * 注：发送到不同的消息队列使用topic标记。
 **/
/**
 * 发送普通消息，用于将消息发送到服务的消息队列
 **/
func SendMsg(c *gin.Context, groupname string, topic string, tag string, data []byte) (*primitive.SendResult, error) {
	// 服务器地址
	endPoint := []string{viper.GetString("rocketmq.addr")}

	// 消息消费失败重试两次
	ProducerMq, err := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(viper.GetInt("rocketmq.retrySize")),
		producer.WithGroupName(groupname),
		producer.WithQueueSelector(producer.NewRandomQueueSelector()),
	)

	defer func(Producer rocketmq.Producer) {
		err := Producer.Shutdown()
		if err != nil {
			fmt.Println("关闭producer失败， error：", err.Error())
			logger.SugarLogger.Error("关闭producer失败，error：", err.Error())
		}
	}(ProducerMq)
	if err != nil {
		fmt.Println("生成producer失败， error：", err.Error())
		logger.SugarLogger.Error("生成producer失败，error：", err.Error())
		return nil, err
	}

	if err := ProducerMq.Start(); err != nil {
		fmt.Println("启动producer失败")
		logger.SugarLogger.Error("启动producer失败，error：", err.Error())
		return nil, err
	}

	msg := primitive.NewMessage(topic, data)
	msg.WithTag(tag)
	res, err := ProducerMq.SendSync(c, msg)
	if err != nil {
		fmt.Println("消息发送失败， error：", err.Error())
		logger.SugarLogger.Error("消息发送失败，error：", err.Error())
		return nil, err
	}

	nowStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s: 消息: %s发送成功 \n", nowStr, res.String())
	return res, nil
}

/**
 * 发送延迟消息
 **/
func SendDelayMsg(topic string, tag string, data []byte) error {
	// 服务器地址
	endPoint := []string{viper.GetString("rocketmq.addr")}

	// 消息消费失败重试两次
	ProducerMq, err := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(viper.GetInt("rocketmq.RetryQueueRetrySize")),
		producer.WithQueueSelector(producer.NewRandomQueueSelector()),
		producer.WithSendMsgTimeout(time.Second*10),
	)

	defer func(Producer rocketmq.Producer) {
		err := Producer.Shutdown()
		if err != nil {
			fmt.Println("关闭producer失败， error：", err.Error())
			logger.SugarLogger.Error("关闭producer失败，error：", err.Error())
		}
	}(ProducerMq)
	if err != nil {
		fmt.Println("生成producer失败， error：", err.Error())
		logger.SugarLogger.Error("生成producer失败，error：", err.Error())
		return err
	}

	if err := ProducerMq.Start(); err != nil {
		fmt.Println("启动producer失败")
		logger.SugarLogger.Error("启动producer失败，error：", err.Error())
		return err
	}

	msg := primitive.NewMessage(topic, data)
	msg.WithTag(tag)
	ProducerMq.SendSync(context.Background(), msg)
	// if err != nil {
	// 	fmt.Println("消息发送失败， error：", err.Error())
	// 	logger.SugarLogger.Error("消息发送失败，error：", err.Error())
	// 	return nil, err
	// }

	// nowStr := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("%s: 消息: %s发送成功 \n", nowStr, res.String())
	// return res, nil
	return nil
}
