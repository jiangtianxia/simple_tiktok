package router

import (
	"fmt"
	"net/http"
	"simple_tiktok/conf"
	"simple_tiktok/controller"
	"simple_tiktok/middlewares"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(cnf *conf.ServerConf) {
	// http
	server := gin.Default()

	// 设置成发布模式
	// gin.SetMode(gin.ReleaseMode)

	// 跨域中间件
	if cnf.OpenCORS {
		server.Use(cors.New(cors.Config{
			AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "PATCH", "OPTIONS", "DELETE"},
			AllowHeaders: []string{"Origin", "Content-Length", "Authorization", "Content-Type", "X-TOKEN",
				"Cookie", "Tus-Extension", "Tus-Resumable", "Tus-Version",
				"Upload-Length", "Upload-Metadata", "Upload-Offset",
				"Access-Control-Allow-Origin", "X-HTTP-Method-Override"},
			AllowCredentials: true,
			AllowAllOrigins:  false,
			AllowOriginFunc: func(origin string) bool {
				for _, once := range cnf.AllowOrigins {
					if once == origin || once == "" {
						return true
					}

					if strings.Contains(origin, once) {
						return true
					}
				}
				return false
			},
			MaxAge: 12 * time.Hour,
		}))
	}

	// 静态资源
	server.StaticFile("/douyin/doc/v2/swagger.json", "apidocs/swagger.json")
	server.StaticFS("/upload", http.Dir("upload"))

	// 路由配置
	v2 := server.Group("/douyin", middlewares.CurrentLimit())
	{
		v2.GET("/hello", controller.Hello)

		// 视频流接口
		v2.GET("/feed", controller.FeedVideo)
		// 用户注册接口
		v2.POST("/user/register/", controller.UserRegister)
		// 用户登录接口
		v2.POST("/user/login/", controller.Userlogin)
		// 用户信息接口
		v2.GET("/user/", controller.GetUserInfo)
		// 视频投稿
		v2.POST("/publish/action/", controller.Publish)
		// 发布列表
		v2.GET("/publish/list/", controller.GetPublishList)

		// 赞操作
		v2.POST("/favorite/action/", controller.Favourite)
		// 喜欢列表
		v2.GET("/favorite/list/", controller.GetFavoriteList)
		// 评论操作
		v2.POST("/comment/action/", controller.CommentAction)
		// 评论列表
		v2.GET("/comment/list/", controller.CommentList)

		// 关注操作
		v2.POST("/relation/action/", controller.UserFollow)
		// 关注列表
		v2.GET("/relation/follow/list/", controller.GetFollowList)
		// 粉丝列表
		v2.GET("/relation/follower/list/", controller.GetFollowerList)
		// 好友列表
		v2.GET("/relation/friend/list/", controller.GetFriendList)
		// 发送消息
		v2.POST("/message/action/", controller.SendMessage)
		// 聊天记录
		v2.GET("/message/chat/", controller.MessageRecord)
	}

	fmt.Printf("http server on port: %d", cnf.HTTPPort)
	server.Run(fmt.Sprintf(":%d", cnf.HTTPPort))
}
