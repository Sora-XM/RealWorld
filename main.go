package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goDemo/config"
	_ "goDemo/docs" // 导入生成的文档包，根据实际包名调整
	"goDemo/route"
	"goDemo/service"
)

// @title RealWorld API
// @version 1.0
// @description RealWorld 后端 API 文档
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /api
func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	// 创建UserService实例，传入数据库连接
	userService := &service.UserService{DB: db}
	profileService := &service.ProfileService{DB: db}

	router := gin.Default()
	// 注册 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	route.SetupRoutes(router, userService)
	route.LoginRoutes(router, userService)
	route.GetCurrentUserRoutes(router, userService)
	route.UpdateUserRoutes(router, userService)
	route.GetProfileRoutes(router, profileService, userService)
	route.FollowUserRoutes(router, profileService, userService)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
