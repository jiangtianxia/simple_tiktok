package rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	redisTopic1 := viper.GetString("rocketmq.redisTopic1")
	redisGroupName1 := viper.GetString("rocketmq.redisGroupName1")
	redisTags1 := viper.GetString("rocketmq.redisTags1")

	go CreateConsumer(redisGroupName, redisTopic, redisTags)
	go CreateConsumer(redisGroupName1, redisTopic1, redisTags1)

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
var loginChan chan ChanMsg = make(chan ChanMsg, 100)

func ReceiveChan() {
	for {
		select {
		case data := <-publishChan:
			// 用户上传视频时，发送videobasic到消息队列，将信息缓存到redis
			videoinfo := &models.VideoBasic{}
			json.Unmarshal(data.Data, videoinfo)
			fmt.Println(videoinfo)
			PublishAction(*videoinfo)

		case data := <-loginChan:
			userLogin := &models.UserBasic{}
			json.Unmarshal(data.Data, userLogin)
			userlogin, err := service.Login(context.Background(), userLogin.Username, userLogin.Password)
			if err != nil {
				context.Background().JSON(http.StatusInternalServerError, gin.H{
					"status_code": -1,
					"status_msg":  err.Error(),
				})
				return
			}

			UserLogin(*userLogin)

			context.Background().JSON(
				http.StatusOK,
				gin.H{
					"status_code": 0,
					"status_msg":  "获取用户信息成功",
					"identity":    userlogin["identity"],
					"token":       userlogin["token"],
				},
			)
		}
	}
}
