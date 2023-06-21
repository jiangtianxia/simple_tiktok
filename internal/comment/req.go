package comment

import "simple_tiktok/internal/common"

type CommentActionReq struct {
	Token       string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	VideoId     string `json:"video_id" form:"video_id" binding:"required"`       // 视频id
	ActionType  int32  `json:"action_type" form:"action_type" binding:"required"` // 1-发布评论, 2-删除评论
	CommentText string `json:"comment_text" form:"comment_text"`                  // optional string comment_text = 4; // 用户填写的评论内容，在action_type=1的时候使用
	CommentId   string `json:"comment_id" form:"comment_id"`                      // 要删除的评论id，在action_type=2的时候使用
}

type CommentListReq struct {
	Token   string `json:"token" form:"token" binding:"required"`       // 用户鉴权token
	VideoId string `json:"video_id" form:"video_id" binding:"required"` // 视频id
	common.SearchReq
	common.PaginateReq
}
