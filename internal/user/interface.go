package user

type IUserService interface {
	// 用户注册
	UserRegister(req *NormalizeUserReq) (*UserRegisterResp, error)

	// 用户登录
	UserLogin(req *NormalizeUserReq) (*UserLoginResp, error)

	// 用户信息
	GetUserInfo(req *UserInfoReq) (*UserInfoResp, error)

	// 用户搜索
	UserSearch(req *UserSearchReq) (*UserSearchResp, error)
}
