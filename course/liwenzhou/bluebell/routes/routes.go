package routes

import (
	"bluebell/controller"
	"bluebell/logger"

	swaggerFiles "github.com/swaggo/files"

	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	// middleware
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	// 注册业务路由
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler) // 注册
	v1.POST("/login", controller.LoginHandler)   // 登录
	//v1.Use(middlewares.JWTAuthMiddleware())      // 中间件
	v1.Use() // 中间件
	{
		// community 社区
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:cid", controller.CommunityDetailHandler)
		// post 帖子
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:pid", controller.PostDetailHandler)
		v1.GET("/posts", controller.GetPostListSimpleHandler)
		v1.GET("/posts2", controller.GetPostListHandler)
		//v1.GET("/posts2/community", controller.GetCommunityPostListHandler)
		// vote 投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	pprof.Register(r) // 注册pprof相关路由

	return r
}
