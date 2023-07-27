package service

import (
	"simple_tiktok/internal/comment"
	"simple_tiktok/internal/favorite"
	"simple_tiktok/internal/user"
	"simple_tiktok/internal/video"
)

var (
	UserService     user.IUserService
	VideoService    video.IVideoService
	FavoriteService favorite.IFavoriteService
	CommentService  comment.ICommentService
)

func init() {
	UserService = &userService{}
	VideoService = &videoService{}
	FavoriteService = &favoriteService{}
	CommentService = &commentService{}
}
