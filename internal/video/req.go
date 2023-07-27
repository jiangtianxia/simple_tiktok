package video

import (
	"mime/multipart"
	"simple_tiktok/internal/common"
)

type VideoFeedReq struct {
	LatestTime int64               `json:"latest_time" form:"latest_time"` // 可选参数, 限制返回视频的最新投稿时间戳, 精确到秒, 不填表示当前时间
	TokenInfo  common.TokenInfoReq `json:"token_info" form:"token_info"`   // tokenInfo
}

type VideoPublishActionReq struct {
	TokenInfo common.TokenInfoReq   `json:"token_info" form:"token_info"`                  // tokenInfo
	File      *multipart.File       `json:"file" form:"file" binding:"required"`           // 视频数据
	FileHead  *multipart.FileHeader `json:"file_head" form:"file_head" binding:"required"` // 视频头
	Title     string                `json:"title" form:"title" binding:"required"`         // 视频标题
}

type VideoPublishListReq struct {
	HashId    string              `json:"hash_id" form:"hash_id" binding:"required"` // 用户hash_id
	UserId    uint                `json:"user_id" form:"user_id"`                    // 用户id
	TokenInfo common.TokenInfoReq `json:"token_info" form:"token_info"`              // tokenInfo
	common.SearchReq
	common.PaginateReq
}

type VideoInfoReq struct {
	HashIds   []string            `json:"hash_ids" form:"hash_ids" binding:"required"` // 视频hash_id
	VideoIds  []uint              `json:"video_ids" form:"video_ids"`                  // 视频id
	TokenInfo common.TokenInfoReq `json:"token_info" form:"token_info"`                // tokenInfo
}
