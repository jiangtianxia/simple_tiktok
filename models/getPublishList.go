package models

// 接收参数
type GetPublishListQuery struct {
	Token string `form:"token"`
	UserId string `form:"user_id"`
}

// 返回参数
type GetPublishListResponse struct {
	StatusCode int `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	VideoList []Video `json:"video_list"`
}

// 视频参数
type Video struct {
	Id int `json:"id"`
	Author Author `json:"author"`
	PlayUrl string `json:"play_url"`
	CoverUrl string `json:"cover_url"`
	FavoriteCount int `json:"favorite_count"`
	CommentCount int `json:"comment_count"`
	IsFavorite bool `json:"is_favorite"`
	Title string `json:"title"`
}

// 作者参数
type Author struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FollowCount int `json:"follow_count"`
	FollowerCount int `json:"follower_count"`
	IsFollow bool `json:"is_follow"`
}
