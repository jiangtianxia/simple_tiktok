package common

type PaginateResp struct {
	Total     int64 `json:"total"`      // 总记录数
	Page      int64 `json:"page"`       // 页数
	PageSize  int64 `json:"page_size"`  // 每页记录数
	TotalPage int64 `json:"total_page"` // 总页数
}

type User struct {
	Id            string `json:"id"`             // 用户id
	Name          string `json:"name"`           // 用户名称
	WorkCount     int64  `json:"work_count"`     // 作品数量
	FavoriteCount int64  `json:"favorite_count"` // 点赞数量
}

type Video struct {
	Id            string `json:"id"`            // 视频唯一标识
	PlayUrl       string `json:"play_url"`      // 视频播放地址
	CoverUrl      string `json:"cover_url"`     // 视频封面地址
	Title         string `json:"title"`         // 视频标题
	FavoriteCount int64  `json:"favorite_count` // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"` // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite"`   // true-已点赞, false-未点赞
	Author        User   `json:"author"`        // 视频作者信息
}

type VideoListResp struct {
	PaginateResp
	VideoList []Video `json:"video_list"` // 视频列表
}
