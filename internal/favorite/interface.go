package favorite

import (
	"simple_tiktok/internal/common"
)

type IFavoriteService interface {
	// 赞操作
	FavoriteAction(req *FavoriteActionReq) error

	// 喜欢列表
	GetFavoriteList(req *FavoriteListReq) (map[string]common.Video, error)
}
