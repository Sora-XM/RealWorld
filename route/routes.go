package route

import (
	"github.com/gin-gonic/gin"
	"goDemo/controller"
	"goDemo/service"
)

// SetupRoutes  注册
func SetupRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.POST("/users/register", userController.RegisterUser)
	}
}

// LoginRoutes 登录
func LoginRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.POST("/users/login", userController.LoginUser)
	}
}

// GetCurrentUserRoutes 获取当前用户
func GetCurrentUserRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.GET("/user", userController.GetCurrentUser)
	}
}

// UpdateUserRoutes 更新用户
func UpdateUserRoutes(router *gin.Engine, UserService *service.UserService) {
	userController := &controller.UserController{UserService: UserService}
	api := router.Group("/api")
	{
		api.PUT("/user", userController.UpdateUser)
	}
}

// GetProfileRoutes 获取用户资料
func GetProfileRoutes(router *gin.Engine, ProfileService *service.ProfileService, UserService *service.UserService) {
	profileController := &controller.ProfileController{ProfileService: ProfileService, UserService: UserService}
	api := router.Group("/api")
	{
		api.GET("/profiles/:username", profileController.GetProfile)
	}
}

// FollowUserRoutes 关注用户
func FollowUserRoutes(router *gin.Engine, ProfileService *service.ProfileService, UserService *service.UserService) {
	profileController := &controller.ProfileController{ProfileService: ProfileService, UserService: UserService}
	api := router.Group("/api")
	{
		api.POST("/follow", profileController.FollowUser)
	}
}
