package user

import "simple_tiktok/internal/common"

type NormalizeUserReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"` // 用户名, 最长32个字符
	Password string `json:"password" form:"password" binding:"required,max=32"` // 登录密码, 最长32个字符
}

type UserInfoReq struct {
	HashIds   []string            `json:"hash_ids" form:"hash_ids" binding:"required"` // 用户hashId
	UserIds   []uint              `json:"user_ids" form:"user_ids"`                    // 用户id
	TokenInfo common.TokenInfoReq `json:"token_info" form:"token_info"`                // tokenInfo
}
