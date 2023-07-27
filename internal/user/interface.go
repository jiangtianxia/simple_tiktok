package user

import "simple_tiktok/internal/common"

type IUserService interface {
	// 用户注册
	UserRegister(req *NormalizeUserReq) (*UserRegisterResp, error)

	// 用户登录
	UserLogin(req *NormalizeUserReq) (*UserLoginResp, error)

	// 用户信息
	GetUserInfo(req *UserInfoReq) (map[string]common.User, error)
}
