package route

import (
	"github.com/gin-gonic/gin"
	"goDemo/controller"
	"goDemo/service"
)

// SetupRoutes  注册路由
func SetupRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.POST("/users/register", userController.RegisterUser)
	}
}

// LoginRoutes 登录路由
func LoginRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.POST("/users/login", userController.LoginUser)
	}
}
