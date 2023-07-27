package grpc_handle

import "simple_tiktok/internal/proto"

var (
	UserGrpcService     proto.UserServiceServer
	VideoGrpcService    proto.VideoServiceServer
	FavoriteGrpcService proto.FavoriteServiceServer
	CommentGrpcService  proto.CommentServiceServer
)

func init() {
	UserGrpcService = &userService{}
	VideoGrpcService = &videoService{}
	FavoriteGrpcService = &favoriteService{}
	CommentGrpcService = &commentService{}
}
