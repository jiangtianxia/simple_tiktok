package handle

import (
	"net/http"
	"path"
	"simple_tiktok/internal/video"
	"simple_tiktok/logger"
	"simple_tiktok/service"
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type videoHandle struct{}

// GetVideoInfo
//
//	@tags			视频相关
//	@Summary		视频信息
//	@Description	视频信息
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string				false	"Authorization"
//	@Param			request			query		video.VideoInfoReq	true	"query params"
//	@Success		200				{object}	map[string]common.Video
//	@Router			/video [get]
func (*videoHandle) GetVideoInfo(ctx *gin.Context) {
	var req video.VideoInfoReq
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

		req.VideoIds = append(req.VideoIds, id)
	}

	resp, err := service.VideoService.GetVideoInfo(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// GetVideoInfo
//
//	@tags			视频相关
//	@Summary		发布列表
//	@Description	发布列表
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string						false	"Authorization"
//	@Param			request			query		video.VideoPublishListReq	true	"query params"
//	@Success		200				{object}	common.VideoListResp
//	@Router			/publish/list [get]
func (*videoHandle) GetVideoPublishList(ctx *gin.Context) {
	var req video.VideoPublishListReq
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

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 30 || req.PageSize <= 0 {
		req.PageSize = 30
	}

	resp, err := service.VideoService.GetVideoPublishList(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// VideoFeed
//
//	@tags			视频相关
//	@Summary		视频流
//	@Description	视频流
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string				false	"Authorization"
//	@Param			request			query		video.VideoFeedReq	true	"query params"
//	@Success		200				{object}	video.VideoFeedResp
//	@Router			/feed [get]
func (*videoHandle) VideoFeed(ctx *gin.Context) {
	var req video.VideoFeedReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	if req.LatestTime == 0 {
		req.LatestTime = time.Now().Unix()
	}

	// 获取token信息
	tokenId, _ := ctx.Get("tokenId")
	tokenUsername, _ := ctx.Get("tokenUsername")
	req.TokenInfo.Id = tokenId.(uint)
	req.TokenInfo.Username = tokenUsername.(string)

	resp, err := service.VideoService.VideoFeed(&req)
	if err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, Result(resp))
}

// VideoPublishAction
//
//	@tags			视频相关
//	@Summary		视频投稿
//	@Description	视频投稿
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization"
//	@Param			data			formData	file	true	"文件"
//	@Param			title			formData	string	true	"标题"
//	@Success		200				string		"success"
//	@version		2.0
//	@Router			/publish/action/ [post]
func (*videoHandle) VideoPublishAction(ctx *gin.Context) {
	title := ctx.DefaultPostForm("title", "")
	req := ctx.Request
	srcFile, head, err := req.FormFile("data")
	if err != nil || title == "" {
		ctx.JSON(http.StatusOK, ErrParam)
		return
	}

	reqData := video.VideoPublishActionReq{
		File:     &srcFile,
		FileHead: head,
		Title:    title,
	}

	// 获取token信息
	tokenId, _ := ctx.Get("tokenId")
	tokenUsername, _ := ctx.Get("tokenUsername")
	reqData.TokenInfo.Id = tokenId.(uint)
	reqData.TokenInfo.Username = tokenUsername.(string)
	if reqData.TokenInfo.Id == 0 || reqData.TokenInfo.Username == "" {
		ctx.AbortWithStatusJSON(http.StatusOK, Result(&MError{
			Code:    -1,
			Message: "请先登录",
			Data:    nil,
		}))
		return
	}

	suffix := path.Ext(head.Filename)
	if suffix != ".avi" && suffix != ".wmv" && suffix != ".mpg" && suffix != ".mpeg" && suffix != ".flv" && suffix != ".mov" && suffix != ".rm" && suffix != ".ram" && suffix != ".swf" && suffix != ".mp4" {
		ctx.AbortWithStatusJSON(http.StatusOK, Result(&MError{
			Code:    -1,
			Message: "请上传视频文件",
			Data:    nil,
		}))
		return
	}

	if err := service.VideoService.VideoPublishAction(&reqData); err != nil {
		logger.SugarLogger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusOK, Result(err))
		return
	}

	ctx.JSON(http.StatusOK, Success)
}
