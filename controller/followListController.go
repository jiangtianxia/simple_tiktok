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
 * @Description 关注列表接口
 * @Date 12:00 2023/2/12
 **/
// 返回结构体
type FollowListRespStruct struct {
	Code     int              `json:"status_code"`
	Msg      string           `json:"status_msg"`
	UserList []service.Author `json:"user_list"`
}

// 传入参数返回
func FollowListResp(c *gin.Context, code int, msg string, userList []service.Author) {
	h := &FollowListRespStruct{
		Code:     code,
		Msg:      msg,
		UserList: userList,
	}

	c.JSON(http.StatusOK, h)
}

// GetFollowList
// @Summary 关注列表
// @Tags 社交接口
// @Param token query string true "token"
// @Param user_id query string true "用户id"
// @Success 200 {object} FollowListRespStruct
// @Router /relation/follow/list/ [get]
func GetFollowList(c *gin.Context) {
	// 1、获取参数
	token := c.DefaultQuery("token", "")
	userId := c.DefaultQuery("user_id", "")

	// 2、验证参数
	if token == "" || userId == "" {
		FollowListResp(c, -1, "请求参数错误", []service.Author{})
		return
	}

	user_id, _ := strconv.Atoi(userId)
	if user_id == 0 {
		FollowListResp(c, -1, "请求参数错误", []service.Author{})
		return
	}

	cnt, err := mysql.FindUserByIdentityCount(uint64(user_id))
	if err != nil {
		logger.SugarLogger.Error("FindUserByIdentityCount Error：", err.Error())
		FollowListResp(c, -1, "请求参数错误", []service.Author{})
		return
	}
	if cnt == 0 {
		FollowListResp(c, -1, "非法用户", []service.Author{})
		return
	}

	// 3、验证token
	// t, _ := utils.GenerateToken(1, "test")
	_, err = middlewares.AuthUserCheck(token)
	if err != nil {
		FollowListResp(c, -1, "无效token", []service.Author{})
		return
	}

	// 4、将数据传入service层
	data, err := service.FollowListService(c, uint64(user_id))
	if err != nil {
		FollowListResp(c, -1, "获取关注列表失败", []service.Author{})
		return
	}

	FollowListResp(c, 0, "获取关注列表成功", data)
}
