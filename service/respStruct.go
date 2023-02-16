package service

/**
 * @Author
 * @Description 视频参数结构体
 * @Date 21:00 2023/2/11
 **/
type VideoInfo struct {
	Id            uint64 `json:"id"`             // 视频唯一标识
	Author        Author `json:"author"`         // 作者信息
	PlayUrl       string `json:"play_url"`       // 视频路径
	CoverUrl      string `json:"cover_url"`      // 封面路径
	FavoriteCount int64  `json:"favorite_count"` // 点赞数
	CommentCount  int64  `json:"comment_count"`  // 评论数
	IsFavorite    bool   `json:"is_favorite"`    // 是否点赞
	Title         string `json:"title"`          // 视频标题
}

/**
 * @Author
 * @Description 评论参数结构体
 * @Date 21:00 2023/2/11
 **/
type CommentInfo struct {
	Id         uint64 `json:"id"`          // 评论唯一标识
	User       Author `json:"user"`        // 用户信息
	Content    string `json:"content"`     // 用户评论内容
	CreateDate string `json:"create_date"` // 评论发布日期
}

/**
 * @Author
 * @Description 作者参数结构体
 * @Date 21:00 2023/2/11
 **/
type Author struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

/**
 * @Author
 * @Description 注册返回参数结构体
 * @Date 21:00 2023/2/11
 **/
type RegisterResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}

/**
 * @Author
 * @Description 登录返回参数结构体
 * @Date 21:00 2023/2/11
 **/
type LoginResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}

/**
 * @Author
 * @Description 聊天记录返回结构体
 * @Date 21:00 2023/2/11
 **/
type Message struct {
	Identity   uint64 `json:"id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
}

/**
 * @Author
 * @Description 好友参数结构体
 * @Date 21:00 2023/2/11
 **/
type Friend struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Avatar        string `json:"avatar"`
	Message       string `json:"message"`
	MsgType       int64  `json:"msgType"`
}
