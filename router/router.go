package router

import (
	"fmt"
	"net/http"
	"simple_tiktok/conf"
	"simple_tiktok/handle"
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
	v2 := server.Group("/douyin/v2", middlewares.CurrentLimit())
	{
		// 用户相关
		{
			// 用户注册
			v2.POST("/user/register", handle.UserHandle.UserRegister)

			// 用户登录
			v2.POST("/user/login", handle.UserHandle.UserLogin)

			// 用户信息
			v2.GET("/user", middlewares.TokenValidator(), handle.UserHandle.GetUserInfo)
		}

		// 视频相关
		{
			// 视频流接口
			v2.GET("/feed", middlewares.TokenValidator(), handle.VideoHandle.VideoFeed)

			// 视频投稿
			v2.POST("/publish/action", middlewares.TokenValidator(), handle.VideoHandle.VideoPublishAction)

			// 发布列表
			v2.GET("/publish/list", middlewares.TokenValidator(), handle.VideoHandle.GetVideoPublishList)

			// 视频信息
			v2.GET("/video", middlewares.TokenValidator(), handle.VideoHandle.GetVideoInfo)
		}

		// 赞相关
		{
			// 赞操作
			v2.POST("/favorite/action", middlewares.TokenValidator(), handle.FavoriteHandle.FavoriteAction)

			// 喜欢列表
			v2.GET("/favorite/list", middlewares.TokenValidator(), handle.FavoriteHandle.FavoriteList)
		}

		// 评论相关
		{
			// 评论操作
			v2.POST("/comment/action", middlewares.TokenValidator(), handle.CommentHandle.CommentAction)

			// 评论列表
			v2.GET("/comment/list", middlewares.TokenValidator(), handle.CommentHandle.CommentList)
		}
	}

	fmt.Printf("http Server started Port:%d", cnf.HTTPPort)
	server.Run(fmt.Sprintf(":%d", cnf.HTTPPort))
}
