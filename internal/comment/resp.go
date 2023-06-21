package comment

import "simple_tiktok/internal/common"

type CommentActionResp struct {
	common.NormalizeResp
	Comment
}

type CommentListResp struct {
	common.NormalizeResp
	common.PaginateResp
	CommentList []Comment `json:"comment_list"` // 评论列表
}

type Comment struct {
	Id         string      `json:"id"`          // 视频评论id
	User       common.User `json:"user"`        // 评论用户信息
	Content    string      `json:"content"`     // 评论内容
	CreateDate string      `json:"create_date"` // 评论发布日期, 格式 mm-dd
}
