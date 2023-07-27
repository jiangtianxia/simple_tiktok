package grpc_handle

import (
	"context"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/favorite"
	"simple_tiktok/internal/proto"
	"simple_tiktok/service"
)

type favoriteService struct {
	proto.UnimplementedFavoriteServiceServer
}

// GetFavoriteList implements proto.FavoriteServiceServer.
func (*favoriteService) GetFavoriteList(ctx context.Context, req *proto.FavoriteListReq) (*proto.VideoInfoResp, error) {
	videoReq := &favorite.FavoriteListReq{
		HashId: "",
		UserId: uint(req.UserId),
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
	}
	videoResp, err := service.FavoriteService.GetFavoriteList(videoReq)
	if err != nil {
		return nil, err
	}

	resp := &proto.VideoInfoResp{}
	for key, video := range videoResp {
		resp.Result[key] = &proto.Video{
			Id:            video.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			Title:         video.Title,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Author: &proto.User{
				Id:            video.Author.Id,
				Name:          video.Author.Name,
				WorkCount:     video.Author.WorkCount,
				FavoriteCount: video.Author.FavoriteCount,
			},
		}
	}
	return resp, nil
}
