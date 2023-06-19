package favorite

import (
	"simple_tiktok/internal/common"
)

type IFavoriteService interface {
	// 赞操作
	FavoriteAction(req *FavoriteActionReq) (*common.NormalizeResp, error)

	// 喜欢列表
	FavoriteList(req *FavoriteListReq) (*common.VideoListResp, error)
}
