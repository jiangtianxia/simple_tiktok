package grpc_handle

import (
	"context"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/proto"
	"simple_tiktok/internal/user"
	"simple_tiktok/service"
)

type userService struct {
	proto.UnimplementedUserServiceServer
}

// GetUserInfo implements proto.UserServiceServer.
func (*userService) GetUserInfo(ctx context.Context, req *proto.UserInfoReq) (*proto.UserInfoResp, error) {
	ids := []uint{}
	for _, id := range req.UserIds {
		ids = append(ids, uint(id))
	}
	userReq := user.UserInfoReq{
		HashIds: []string{},
		UserIds: ids,
		TokenInfo: common.TokenInfoReq{
			Id:       uint(req.TokenInfo.Id),
			Username: req.TokenInfo.Username,
		},
	}

	userResp, err := service.UserService.GetUserInfo(&userReq)
	if err != nil {
		return nil, err
	}

	resp := &proto.UserInfoResp{}
	for key, user := range userResp {
		resp.Result[key] = &proto.User{
			Id:            user.Id,
			Name:          user.Name,
			WorkCount:     user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		}
	}
	return resp, nil
}
