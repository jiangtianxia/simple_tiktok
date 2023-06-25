package user

import "simple_tiktok/internal/common"

type UserRegisterResp struct {
	common.NormalizeResp
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginResp struct {
	common.NormalizeResp
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResp struct {
	common.NormalizeResp
	common.PaginateResp
	UserList []common.User `json:"user_list"`
}
