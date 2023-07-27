package grpc_handle

import (
	"context"
	"simple_tiktok/internal/comment"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/proto"
	"simple_tiktok/service"
)

type commentService struct {
	proto.UnimplementedCommentServiceServer
}

// GetCommentList implements proto.CommentServiceServer.
func (*commentService) GetCommentList(ctx context.Context, req *proto.CommentListReq) (*proto.CommentListResp, error) {
	commentReq := comment.CommentListReq{
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
		HashId:  "",
		VideoId: uint(req.VideoId),
		PaginateReq: common.PaginateReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	commentResp, err := service.CommentService.GetCommentList(&commentReq)
	if err != nil {
		return nil, err
	}

	var commentList []*proto.Comment
	for _, comment := range commentResp.CommentList {
		commentList = append(commentList, &proto.Comment{
			Id: comment.Id,
			User: &proto.User{
				Id:            comment.User.Id,
				Name:          comment.User.Name,
				WorkCount:     comment.User.WorkCount,
				FavoriteCount: comment.User.FavoriteCount,
			},
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		})
	}

	return &proto.CommentListResp{
		Total:       commentResp.Total,
		Page:        commentResp.Page,
		PageSize:    commentResp.PageSize,
		TotalPage:   commentResp.TotalPage,
		CommentList: commentList,
	}, nil
}
