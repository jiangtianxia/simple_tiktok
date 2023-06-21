package favorite

import "simple_tiktok/internal/common"

type FavoriteActionReq struct {
	Token      string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	VideoId    int64  `json:"video_id" form:"video_id" binding:"required"`       // 视频id
	ActionType int32  `json:"action_type" form:"action_type" binding:"required"` // 1-点赞, 2-取消点赞
}

type FavoriteListReq struct {
	UserId int64  `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token  string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
	common.SearchReq
	common.PaginateReq
}
