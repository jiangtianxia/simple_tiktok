package user

import "simple_tiktok/internal/common"

type NormalizeUserReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"` // 用户名, 最长32个字符
	Password string `json:"password" form:"password" binding:"required,max=32"` // 登录密码, 最长32个字符
}

type UserInfoReq struct {
	UserId string `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token  string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
}

type UserSearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"` // 搜索关键词, 用户名/用户id
	common.SearchReq
	common.PaginateReq
}
