package service

/**
 * @Author
 * @Description 注册请求体
 * @Date 21:00 2023/2/11
 **/
type RegisterRequire struct {
	Username string
	Password string
}

/**
 * @Author jiang
 * @Description 点赞请求体
 * @Date 21:00 2023/2/11
 **/
type FollowReqStruct struct {
	UserId     string
	ToUserId   string
	ActionType int
}
