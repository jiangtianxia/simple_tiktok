package controller

import (
	"net/http"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

// 返回结构体
type UserInfoRespStruct struct {
	Code int            `json:"status_code"`
	Msg  string         `json:"status_msg"`
	User service.Author `json:"user"`
}

// 传入参数返回
func UserInfoResp(c *gin.Context, code int, msg string, user service.Author) {
	h := &UserInfoRespStruct{
		Code: code,
		Msg:  msg,
		User: user,
	}

	c.JSON(http.StatusOK, h)
}

/**
 * @Author Xiaoyu Zhang
 * @Description 用户信息接口
 * @Date 14:00 2023/1/31
 **/
// GetUserInfo
// @Summary 用户信息
// @Tags 基础接口
// @Param token query string true "token"
// @Param user_id query string true "用户id"
// @Success 200 {object} UserInfoRespStruct
// @Router /user/ [get]
func GetUserInfo(c *gin.Context) {
	// 接受参数
	userId := c.DefaultQuery("user_id", "0")
	token := c.DefaultQuery("token", "")

	// 验证token
	userClaim, err := middlewares.AuthUserCheck(token)
	if err != nil {
		UserInfoResp(c, -1, err.Error(), service.Author{})
		return
	}

	// 把参数传给service层
	res, err := service.UserInfo(c, userClaim.Identity, userId)
	if err != nil {
		UserInfoResp(c, -1, "获取用户信息失败", service.Author{})
		return
	}

	UserInfoResp(c, 0, "获取用户信息成功", res)
}
