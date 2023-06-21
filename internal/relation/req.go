package relation

import "simple_tiktok/internal/common"

type RelationActionReq struct {
	Token      string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	ToUserId   string `json:"to_user_id" form:"to_user_id" binding:"required"`   // 对方用户id
	ActionType int32  `json:"action_type" form:"action_type" binding:"required"` // 1-关注, 2-取消关注
}

type RelationReq struct {
	UserId string `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token  string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
	common.SearchReq
	common.PaginateReq
}
