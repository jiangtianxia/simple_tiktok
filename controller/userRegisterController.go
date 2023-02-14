package controller

import (
	"net/http"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

/**
 * Creator: lx
 * Last Editor: lx
 * Description: controller层，解析参数，处理参数，并打包传给service层
 **/
// 返回结构体
type UserRegisterRespStruct struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId uint64 `json:"user_id"`
	Token  string `json:"token"`
}

// 传入参数返回
func UserRegisterResp(c *gin.Context, code int, msg string, userId uint64, token string) {
	h := &UserRegisterRespStruct{
		Code:   code,
		Msg:    msg,
		UserId: userId,
		Token:  token,
	}

	c.JSON(http.StatusOK, h)
}

// UserRegister
// @Summary 用户注册
// @Tags 基础接口
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 200 {string} status_code status_msg
// @Router /user/register/ [post]
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	req := service.RegisterRequire{
		Username: string(username),
		Password: string(password),
	}

	registerResponse, err := service.PostUserRegister(c, &req)
	if err != nil {
		UserRegisterResp(c, -1, err.Error(), 0, "")
		return
	}

	UserRegisterResp(c, 0, "注册成功", registerResponse.Identity, registerResponse.Token)
}
