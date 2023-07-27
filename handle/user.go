package handle

import (
	"fmt"
	"net/http"
	"simple_tiktok/internal/user"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

type userHandle struct{}

// UserRegister
//
//	@tags			用户相关
//	@Summary		用户注册
//	@Description	用户注册
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			request	body		user.NormalizeUserReq	true	"body"
//	@Success		200		{object}	user.UserRegisterResp
//	@Router			/user/register [post]
func (*userHandle) UserRegister(ctx *gin.Context) {
	var req user.NormalizeUserReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	resp, err := service.UserService.UserRegister(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// UserRegister
//
//	@tags			用户相关
//	@Summary		用户登录
//	@Description	用户登录
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			request	body		user.NormalizeUserReq	true	"body"
//	@Success		200		{object}	user.UserLoginResp
//	@Router			/user/login [post]
func (*userHandle) UserLogin(ctx *gin.Context) {
	var req user.NormalizeUserReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	// 判断是否已经登录限制
	flag, time := utils.ExceLoginBank(ctx, req.Username)
	if flag {
		ctx.AbortWithStatusJSON(http.StatusOK, Result(MError{
			Code:    -1,
			Message: fmt.Sprintf("1分钟内连续登录5次, 登录限制, 请 %d 分钟后重试", int(time)),
			Data:    nil,
		}))
		return
	}

	resp, err := service.UserService.UserLogin(&req)
	if err != nil {
		utils.AddLoginError(ctx, req.Username)
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// GetUserInfo
//
//	@tags			用户相关
//	@Summary		用户信息
//	@Description	用户信息
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string				false	"Authorization"
//	@Param			request			query		user.UserInfoReq	true	"query"
//	@Success		200				{object}	map[string]common.User
//	@Router			/user [get]
func (*userHandle) GetUserInfo(ctx *gin.Context) {
	var req user.UserInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	// 获取token信息
	tokenId, _ := ctx.Get("tokenId")
	tokenUsername, _ := ctx.Get("tokenUsername")
	req.TokenInfo.Id = tokenId.(uint)
	req.TokenInfo.Username = tokenUsername.(string)

	// 解析hashId
	for _, hashId := range req.HashIds {
		id, err := utils.DecodeID(hashId)
		if err != nil {
			ctx.JSON(http.StatusOK, ErrParam)
			return
		}

		req.UserIds = append(req.UserIds, id)
	}

	resp, err := service.UserService.GetUserInfo(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}
