package router

import (
	"net/http"

	"simple_tiktok/controller"
	docs "simple_tiktok/docs"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// swagger 配置
	docs.SwaggerInfo.BasePath = "/douyin"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//加载静态资源，一般是上传的资源，例如用户上传的图片
	r.StaticFS("/upload", http.Dir("upload"))

	// 路由配置
	v1 := r.Group("/douyin")
	{
		/*
		* 公共接口
		 */
		v1.GET("/hello", controller.Hello)

	}

	r.GET("/test", controller.Test)

	return r
}
