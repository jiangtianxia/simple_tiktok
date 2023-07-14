package controller

import (
	"fmt"
	"net/http"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 视频投稿接口
 * @Date 16:00 2023/1/28
 **/
// 返回结构体
type UploadRespStruct struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// 传入参数返回
func UploadResp(c *gin.Context, code int, msg string) {
	h := &UploadRespStruct{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, h)
}

// Publish
//	@Summary	视频投稿
//	@Tags		基础接口
//	@Param		token	formData	string	true	"token"
//	@Param		data	formData	file	true	"文件"
//	@Param		title	formData	string	true	"标题"
//	@Success	200		{object}	UploadRespStruct
//	@Router		/publish/action/ [post]
func Publish(c *gin.Context) {
	// 1、获取参数
	token := c.DefaultPostForm("token", "")
	title := c.PostForm("title")

	req := c.Request
	srcFile, head, err := req.FormFile("data")
	if err != nil {
		logger.SugarLogger.Error("Get FormFile Error:" + err.Error())
		fmt.Println("Get FormFile Error:" + err.Error())
		UploadResp(c, -1, "请求参数错误")
		return
	}

	// 2、验证token
	if token == "" {
		UploadResp(c, -1, "无效的Token")
		return
	}

	userClaim, err := middlewares.AuthUserCheck(token)
	if err != nil {
		UploadResp(c, -1, err.Error())
		return
	}

	// 3、将数据传到service层
	code, msg := service.UploadCOS(c, srcFile, head, title, userClaim.Identity)
	UploadResp(c, code, msg)
}
