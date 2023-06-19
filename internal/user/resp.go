package user

import "simple_tiktok/internal/common"

type UserRegisterResp struct {
	common.NormalizeResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginResp struct {
	common.NormalizeResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResp struct {
	common.NormalizeResp
	common.User
}

type UserSearchResp struct {
	common.NormalizeResp
	common.PaginateResp
	Users []common.User `json:"users"`
}
