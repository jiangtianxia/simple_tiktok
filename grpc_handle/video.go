package grpc_handle

import (
	"context"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/proto"
	"simple_tiktok/internal/video"
	"simple_tiktok/service"
)

type videoService struct {
	proto.UnimplementedVideoServiceServer
}

// GetVideoInfo implements proto.VideoServiceServer.
func (*videoService) GetVideoInfo(ctx context.Context, req *proto.VideoInfoReq) (*proto.VideoInfoResp, error) {
	ids := []uint{}
	for _, id := range req.VideoIds {
		ids = append(ids, uint(id))
	}
	videoReq := &video.VideoInfoReq{
		HashIds:  []string{},
		VideoIds: ids,
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
	}
	videoResp, err := service.VideoService.GetVideoInfo(videoReq)
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

// GetVideoPublishList implements proto.VideoServiceServer.
func (*videoService) GetVideoPublishList(ctx context.Context, req *proto.VideoPublishListReq) (*proto.VideoListResp, error) {
	if req.Page <= 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 30 {
		req.PageSize = 30
	}

	videoReq := &video.VideoPublishListReq{
		HashId: "",
		UserId: uint(req.UserId),
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
		SearchReq: common.SearchReq{
			Where: req.Where,
			Order: req.Order,
			Sort:  int(req.Sort),
		},
		PaginateReq: common.PaginateReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	videoResp, err := service.VideoService.GetVideoPublishList(videoReq)
	if err != nil {
		return nil, err
	}

	var videoList []*proto.Video
	for _, v := range videoResp.VideoList {
		videoList = append(videoList, &proto.Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Author: &proto.User{
				Id:            v.Author.Id,
				Name:          v.Author.Name,
				WorkCount:     v.Author.WorkCount,
				FavoriteCount: v.Author.FavoriteCount,
			},
		})
	}
	return &proto.VideoListResp{
		Total:     videoResp.Total,
		Page:      videoResp.Page,
		PageSize:  videoResp.PageSize,
		TotalPage: videoResp.TotalPage,
		VideoList: videoList,
	}, nil
}

// VideoFeed implements proto.VideoServiceServer.
func (*videoService) VideoFeed(ctx context.Context, req *proto.VideoFeedReq) (*proto.VideoFeedResp, error) {
	videoReq := video.VideoFeedReq{
		LatestTime: req.LatestTime,
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
	}
	videoResp, err := service.VideoService.VideoFeed(&videoReq)
	if err != nil {
		return nil, err
	}

	var videoList []*proto.Video
	for _, v := range videoResp.VideoList {
		videoList = append(videoList, &proto.Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Author: &proto.User{
				Id:            v.Author.Id,
				Name:          v.Author.Name,
				WorkCount:     v.Author.WorkCount,
				FavoriteCount: v.Author.FavoriteCount,
			},
		})
	}
	return &proto.VideoFeedResp{
		NextTime:  videoResp.NextTime,
		VideoList: videoList,
	}, nil
}
