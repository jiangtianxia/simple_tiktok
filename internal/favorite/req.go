package favorite

import "simple_tiktok/internal/common"

type FavoriteActionReq struct {
	TokenInfo  common.TokenInfoReq `json:"token_info" form:"token_info"`                      // tokenInfo
	HashId     string              `json:"hash_id" form:"hash_id" binding:"required"`         // 视频hashId
	VideoId    uint                `json:"video_id" form:"video_id"`                          // 视频id
	ActionType int32               `json:"action_type" form:"action_type" binding:"required"` // 1-点赞, 2-取消点赞
}

type FavoriteListReq struct {
	HashId    string              `json:"hash_id" form:"hash_id" binding:"required"` // 用户hashId
	UserId    uint                `json:"user_id" form:"user_id"`                    // 用户id
	TokenInfo common.TokenInfoReq `json:"token_info" form:"token_info"`              // tokenInfo
}
