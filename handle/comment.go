package handle

import (
	"net/http"
	"simple_tiktok/internal/comment"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"simple_tiktok/utils"

	"github.com/gin-gonic/gin"
)

type commentHandle struct {
}

// CommentAction
//
//	@tags			评论相关
//	@Summary		评论操作
//	@Description	评论操作
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string						true	"Authorization"
//	@Param			request			body		comment.CommentActionReq	true	"body"
//	@Success		200				{object}	comment.Comment
//	@Router			/comment/action [post]
func (*commentHandle) CommentAction(ctx *gin.Context) {
	var req comment.CommentActionReq
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

	if req.ActionType == 1 && req.CommentText == "" {
		ctx.AbortWithStatusJSON(http.StatusOK, Result(MError{
			Code:    -1,
			Message: "评论内容为空",
			Data:    nil,
		}))
		return
	}
	if req.ActionType == 2 && req.CommentId == "" {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	resp, err := service.CommentService.CommentAction(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// CommentList
//
//	@tags			评论相关
//	@Summary		评论列表
//	@Description	评论列表
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string					false	"Authorization"
//	@Param			request			query		comment.CommentListReq	true	"query"
//	@Success		200				{object}	comment.CommentListResp
//	@Router			/comment/list [get]
func (*commentHandle) CommentList(ctx *gin.Context) {
	var req comment.CommentListReq
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
	req.VideoId = id

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 30 || req.PageSize <= 0 {
		req.PageSize = 30
	}

	resp, err := service.CommentService.CommentList(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}
