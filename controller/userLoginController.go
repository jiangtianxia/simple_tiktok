package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/service"
)

type LoginRespStruct struct {
	Code string
	Msg  string
}

func LoginResp(c *gin.Context, code string, msg string) {
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
	LoginResp(c, code, msg)

}
