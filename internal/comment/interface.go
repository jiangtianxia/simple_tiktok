package comment

type ICommentService interface {
	// 视频评论操作
	CommentAction(req *CommentActionReq) (*Comment, error)

	// 视频评论列表
	GetCommentList(req *CommentListReq) (*CommentListResp, error)
}
