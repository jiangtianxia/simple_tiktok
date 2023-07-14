package controller

import (
	"net/http"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/middlewares"
	"simple_tiktok/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 好友列表接口
 * @Date 17:00 2023/2/12
 **/
// 返回结构体
type FirendListRespStruct struct {
	Code     int              `json:"status_code"`
	Msg      string           `json:"status_msg"`
	UserList []service.Friend `json:"user_list"`
}

// 传入参数返回
func FirendListResp(c *gin.Context, code int, msg string, userList []service.Friend) {
	h := &FirendListRespStruct{
		Code:     code,
		Msg:      msg,
		UserList: userList,
	}

	c.JSON(http.StatusOK, h)
}

// GetFriendList
//	@Summary	好友列表
//	@Tags		社交接口
//	@Param		token	query		string	true	"token"
//	@Param		user_id	query		string	true	"用户id"
//	@Success	200		{object}	FirendListRespStruct
//	@Router		/relation/friend/list/ [get]
func GetFriendList(c *gin.Context) {
	// 1、获取参数
	token := c.DefaultQuery("token", "")
	userId := c.DefaultQuery("user_id", "")

	// 2、验证参数
	if token == "" || userId == "" {
		FirendListResp(c, -1, "请求参数错误", []service.Friend{})
		return
	}

	user_id, _ := strconv.Atoi(userId)
	if user_id == 0 {
		FirendListResp(c, -1, "请求参数错误", []service.Friend{})
		return
	}

	cnt, err := mysql.FindUserByIdentityCount(uint64(user_id))
	if err != nil {
		logger.SugarLogger.Error("FindUserByIdentityCount Error：", err.Error())
		FirendListResp(c, -1, "请求参数错误", []service.Friend{})
		return
	}
	if cnt == 0 {
		FirendListResp(c, -1, "非法用户", []service.Friend{})
		return
	}

	// 3、验证token
	// t, _ := utils.GenerateToken(1, "test")
	_, err = middlewares.AuthUserCheck(token)
	if err != nil {
		FirendListResp(c, -1, "无效token", []service.Friend{})
		return
	}

	// 4、将数据传入service层
	data, err := service.FriendListService(c, uint64(user_id))
	if err != nil {
		FirendListResp(c, -1, "获取好友列表失败", []service.Friend{})
		return
	}

	FirendListResp(c, 0, "获取好友列表成功", data)
}
