package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/models"
	rocket "simple_tiktok/rocketmq"
	"time"
)

func Userlogin(c *gin.Context) {
	//获取参数
	username := c.Query("username")
	password := c.Query("password")

	user := models.UserBasic{
		Password: password,
		Username: username,
	}

	data, _ := json.Marshal(user)

	// 发送消息
	res, err := rocket.SendMsg(c, "LoginProducer", "LoginTopic", "login", data)
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
