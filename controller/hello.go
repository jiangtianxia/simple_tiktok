package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	rocket "simple_tiktok/rocketmq"
	"time"

	"github.com/gin-gonic/gin"
)

// Hello
// @Tags 公共接口
// @Summary 首页
// @Success 200 {string} hello world
// @Router /hello [get]
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

type Name struct {
	Msgid    string `json:"msgid"`
	Userid   string `json:"userid"`
	Username string `json:"username"`
}

func Test(c *gin.Context) {
	user := &Name{
		Userid:   "20230123",
		Username: "jiang",
	}

	data, _ := json.Marshal(user)

	// topic := viper.GetString("rocketmq.ServerTopic")
	// 发送消息
	res, err := rocket.SendMsg(c, "ServerProducer", "SimpleTopic", "test", data)
	if err != nil {
		fmt.Println("发送消息失败， error:", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "失败",
		})
		return
	}

	fmt.Println(res.Status)
	fmt.Println(res.Status == 1)
	if res.Status == 0 {
		for i := 0; i < 30; i++ {
			// hash msgid req
			// 根据msgid，查询redis缓存中是否存在数据，如果存在则将结果返回
			c.JSON(http.StatusOK, gin.H{
				"message": "成功",
				"res":     res.MsgID,
			})
			fmt.Println("查询缓存")
			time.Sleep(time.Second)
		}

		// 30秒后，还是没有结果，则返回请求超时
		c.JSON(http.StatusOK, gin.H{
			"message": "请求超时",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送失败",
	})
}
