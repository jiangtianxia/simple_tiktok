package rocket

import (
	"encoding/json"
	"fmt"

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
	serverTopic := viper.GetString("rocketmq.serverTopic")
	serverGroupName := viper.GetString("rocketmq.serverGroupName")
	serverTags := viper.GetString("rocketmq.serverTags")
	redisTopic := viper.GetString("rocketmq.redisTopic")
	redisGroupName := viper.GetString("rocketmq.redisGroupName")
	redisTags := viper.GetString("rocketmq.redisTags")

	go CreateConsumer(serverGroupName, serverTopic, serverTags)
	go CreateConsumer(redisGroupName, redisTopic, redisTags)
	go CreateConsumer("Test", "SimpleTopic", "test||login")
	go CreateConsumer("test", "SimpleTopic", "test||login")
	go CreateConsumer("test1", "SimpleTopic", "test||login")
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

var sendChan chan ChanMsg = make(chan ChanMsg, 100)
var LoginChan chan ChanMsg = make(chan ChanMsg, 100)

func ReceiveChan() {
	for {
		select {
		case data := <-sendChan:
			info := &Name{}
			json.Unmarshal(data.Data, info)
			fmt.Println(info)
			fmt.Println("sendChan")
		case data := <-LoginChan:
			info := &Name{}
			json.Unmarshal(data.Data, info)
			fmt.Println(info)
			fmt.Println("LoginChan")
		}
	}
}

type Name struct {
	Userid   string `json:"userid"`
	Username string `json:"username"`
}
