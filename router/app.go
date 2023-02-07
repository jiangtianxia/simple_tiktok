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
		v1.GET("/user", controller.UserInfo)
		//用户操作
		v1.POST("/user/login", controller.Userlogin)

		// 视频投稿
		v1.POST("/publish/action/", controller.Publish)
		//赞操作
		v1.POST("/favorite/action/", controller.Favourite)
	}

	r.GET("/test", controller.Test)

	return r
}
