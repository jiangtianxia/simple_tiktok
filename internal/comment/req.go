package comment

import "simple_tiktok/internal/common"

type CommentActionReq struct {
	TokenInfo   common.TokenInfoReq `json:"token_info" form:"token_info"`                      // tokenInfo
	HashId      string              `json:"hash_id" form:"hash_id" binding:"required"`         // 视频hashId
	VideoId     uint                `json:"video_id" form:"video_id"`                          // 视频id
	ActionType  int32               `json:"action_type" form:"action_type" binding:"required"` // 1-发布评论, 2-删除评论
	CommentText string              `json:"comment_text" form:"comment_text"`                  // optional string comment_text = 4; // 用户填写的评论内容，在action_type=1的时候使用
	CommentId   string              `json:"comment_id" form:"comment_id"`                      // 要删除的评论id，在action_type=2的时候使用
}

type CommentListReq struct {
	TokenInfo common.TokenInfoReq `json:"token_info" form:"token_info"`              // tokenInfo
	HashId    string              `json:"hash_id" form:"hash_id" binding:"required"` // 视频hashId
	VideoId   uint                `json:"video_id" form:"video_id"`                  // 视频id
	common.PaginateReq
}
