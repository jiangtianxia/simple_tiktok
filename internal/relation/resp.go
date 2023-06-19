package relation

import "simple_tiktok/internal/common"

type RelationResp struct {
	common.NormalizeResp
	common.PaginateResp
	UserList common.User `json:"user_list"`
}
