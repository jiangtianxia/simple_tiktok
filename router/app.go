package router

import (
	"simple_tiktok/controller"
	docs "simple_tiktok/docs"
	"simple_tiktok/middlewares"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 设置成发布模式
	// gin.SetMode(gin.ReleaseMode)

	// 全局使用熔断器，加入熔断保障
	r.Use(middlewares.GinCircuitBreaker)

	// swagger 配置
	docs.SwaggerInfo.BasePath = "/douyin"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 路由配置
	v1 := r.Group("/douyin", middlewares.CurrentLimit())
	{
		/*
		* 公共接口
		 */
		v1.GET("/hello", controller.Hello)

		/*
		* 基础接口
		 */
		// 视频流接口
		v1.GET("/feed", controller.FeedVideo)

		// 用户注册接口
		v1.POST("/user/register", controller.UserRegister)

		// 用户登录接口
		v1.POST("/user/login", controller.Userlogin)

		// 用户信息接口
		v1.GET("/user", controller.GetUserInfo)

		// 视频投稿
		v1.POST("/publish/action/", controller.Publish)

		// 发布列表
		v1.GET("/publish/list/", controller.GetPublishList)

		/*
		* 社交接口
		 */
		// 关注操作
		v1.POST("/relation/action/", controller.UserFollow)

		// 关注列表
		v1.GET("/relation/follow/list/", controller.GetFollowList)

		// 粉丝列表
		v1.GET("/relation/follower/list/", controller.GetFollowerList)

		// 好友列表
		v1.GET("/relation/friend/list/", controller.GetFriendList)
	}
	return r
}
