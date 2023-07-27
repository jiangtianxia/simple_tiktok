package handle

import "github.com/gin-gonic/gin"

var (
	UserHandle     IUserHandle
	VideoHandle    IVideoHandle
	FavoriteHandle IFavoriteHandle
	CommentHandle  ICommentHandle
)

func init() {
	UserHandle = &userHandle{}
	VideoHandle = &videoHandle{}
	FavoriteHandle = &favoriteHandle{}
	CommentHandle = &commentHandle{}
}

// 用户相关
type IUserHandle interface {
	// 用户注册
	UserRegister(ctx *gin.Context)

	// 用户登录
	UserLogin(ctx *gin.Context)

	// 用户信息
	GetUserInfo(ctx *gin.Context)
}

// 视频相关
type IVideoHandle interface {
	// 视频流
	VideoFeed(ctx *gin.Context)

	// 视频投稿
	VideoPublishAction(ctx *gin.Context)

	// 发布列表
	GetVideoPublishList(ctx *gin.Context)

	// 视频信息
	GetVideoInfo(ctx *gin.Context)
}

// 赞相关
type IFavoriteHandle interface {
	// 赞操作
	FavoriteAction(ctx *gin.Context)

	// 喜欢列表
	FavoriteList(ctx *gin.Context)
}

// 评论相关
type ICommentHandle interface {
	// 评论操作
	CommentAction(ctx *gin.Context)

	// 评论列表
	CommentList(ctx *gin.Context)
}
