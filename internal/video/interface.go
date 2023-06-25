package video

import "simple_tiktok/internal/common"

type IVideoService interface {
	// 视频流
	VideoFeed(req *VideoFeedReq) (*VideoFeedResp, error)

	// 视频投稿
	VideoPublishAction(req *VideoPublishActionReq) (*common.NormalizeResp, error)

	// 发布列表
	GetVideoPublishList(req *VideoPublishListReq) (*common.VideoListResp, error)

	// 视频信息
	GetVideoInfo(req *VideoInfoReq) (*common.VideoListResp, error)

	// 视频搜索
	VideoSearch(req *VideoSearchReq) (*common.VideoListResp, error)
}
