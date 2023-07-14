package controller

import (
	"net/http"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

// 返回结构体
type UserLoginRespStruct struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId uint64 `json:"user_id"`
	Token  string `json:"token"`
}

// 传入参数返回
func UserLoginResp(c *gin.Context, code int, msg string, userId uint64, token string) {
	h := &UserLoginRespStruct{
		Code:   code,
		Msg:    msg,
		UserId: userId,
		Token:  token,
	}

	c.JSON(http.StatusOK, h)
}

// Userlogin
//	@Summary	用户登录
//	@Tags		基础接口
//	@Param		username	query		string	true	"username"
//	@Param		password	query		string	true	"password"
//	@Success	200			{object}	UserLoginRespStruct
//	@Router		/user/login/ [post]
func Userlogin(c *gin.Context) {
	//获取参数
	username := c.Query("username")
	password := c.Query("password")

	data, err := service.Login(c, username, password)
	if err != nil {
		UserLoginResp(c, -1, err.Error(), 0, "")
		return
	}

	UserLoginResp(c, 0, "登录成功", data.Identity, data.Token)
}
