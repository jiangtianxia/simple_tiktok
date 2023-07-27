package handle

import (
	"net/http"
	"simple_tiktok/internal/favorite"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

type favoriteHandle struct{}

// FavoriteAction
//
//	@tags			赞相关
//	@Summary		赞操作
//	@Description	赞操作
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header	string						true	"Authorization"
//	@Param			request			body	favorite.FavoriteActionReq	true	"body"
//	@Success		200				string	success
//	@Router			/favorite/action [post]
func (*favoriteHandle) FavoriteAction(ctx *gin.Context) {
	var req favorite.FavoriteActionReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	// 获取token信息
	tokenId, _ := ctx.Get("tokenId")
	tokenUsername, _ := ctx.Get("tokenUsername")
	req.TokenInfo.Id = tokenId.(uint)
	req.TokenInfo.Username = tokenUsername.(string)

	if req.TokenInfo.Id == 0 || req.TokenInfo.Username == "" {
		ctx.AbortWithStatusJSON(http.StatusOK, Result(MError{
			Code:    -1,
			Message: "请先登录",
			Data:    nil,
		}))
		return
	}
	// 解析hashId
	id, err := utils.DecodeID(req.HashId)
	if err != nil || (req.ActionType != 1 && req.ActionType != 2) {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}
	req.VideoId = id

	if err := service.FavoriteService.FavoriteAction(&req); err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.JSON(http.StatusOK, Success)
}

// FavoriteList
//
//	@tags			赞相关
//	@Summary		点赞列表
//	@Description	点赞列表
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string						false	"Authorization"
//	@Param			request			query		favorite.FavoriteListReq	true	"query params"
//	@Success		200				{object}	map[string]common.Video
//	@Router			/favorite/list [get]
func (*favoriteHandle) FavoriteList(ctx *gin.Context) {
	var req favorite.FavoriteListReq
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
	id, err := utils.DecodeID(req.HashId)
	if err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}
	req.UserId = id

	resp, err := service.FavoriteService.GetFavoriteList(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}
