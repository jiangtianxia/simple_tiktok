package video

import "simple_tiktok/internal/common"

type VideoFeedReq struct {
	LatestTime int64  `json:"latest_time" form:"latest_time"` // 可选参数, 限制返回视频的最新投稿时间戳, 精确到秒, 不填表示当前时间
	Token      string `json:"token" form:"token"`             // 用户鉴权token
}

type VideoPublishActionReq struct {
	Token string `json:"token" form:"token" binding:"required"` // 用户鉴权token
	Data  []byte `json:"data" form:"data" binding:"required"`   // 视频数据
	Title string `json:"title" form:"title" binding:"required"` // 视频标题
}

type VideoPublishListReq struct {
	UserId int64  `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token  string `json:"token" form:"token"`                        // 用户鉴权token
	common.SearchReq
	common.PaginateReq
}

type VideoInfoReq struct {
	VideoId int64  `json:"video_id" form:"video_id" binding:"required"` // 视频id
	Token   string `json:"token" form:"token" binding:"required"`       // 用户鉴权token
}

type VideoSearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"` // 搜索关键词, 视频标题/作者/视频id
	common.SearchReq
	common.PaginateReq
}
