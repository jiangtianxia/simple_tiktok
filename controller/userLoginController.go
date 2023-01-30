package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/service"
	"simple_tiktok/utils"
)

type LoginRespStruct struct {
	Code int64
	Msg  string
}

func LoginResp(c *gin.Context, code int64, msg string) {
	h := &LoginRespStruct{
		Code: code,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, h)
}

func Userlogin(c *gin.Context) {
	//获取参数
	username := c.Query("username")
	getpassword, _ := c.Get("password")
	password, ok := getpassword.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码格式错误",
		})
		return
	}
	code, msg := service.Login(c, username, password)

	//布隆过滤器
	if !utils.Filter.Check(username) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  "no such user",
		})
		return
	}

	LoginResp(c, code, msg)

}
