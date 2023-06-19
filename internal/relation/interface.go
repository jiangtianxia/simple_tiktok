package relation

import "simple_tiktok/internal/common"

type IRelationService interface {
	// 关系操作
	RelationAction(req *RelationActionReq) (*common.NormalizeResp, error)

	// 用户关注列表
	RelationFollowList(req *RelationReq) (*RelationResp, error)

	// 用户粉丝列表
	RelationFollowerList(req *RelationReq) (*RelationResp, error)
}
