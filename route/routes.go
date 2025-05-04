package route

import (
	"github.com/gin-gonic/gin"
	"goDemo/controller"
	"goDemo/service"
	"goDemo/utils"
)

// SetupRoutes  注册
func SetupRoutes(router *gin.Engine, UserService *service.UserService, Auth *utils.Auth) {
	userController := &controller.UserController{UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.POST("/users/register", userController.RegisterUser)
	}
}

// LoginRoutes 登录
func LoginRoutes(router *gin.Engine, UserService *service.UserService, Auth *utils.Auth) {
	userController := &controller.UserController{UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.POST("/users/login", userController.LoginUser)
	}
}

// GetCurrentUserRoutes 获取当前用户
func GetCurrentUserRoutes(router *gin.Engine, UserService *service.UserService, Auth *utils.Auth) {
	userController := &controller.UserController{UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.GET("/user", userController.GetCurrentUser)
	}
}

// UpdateUserRoutes 更新用户
func UpdateUserRoutes(router *gin.Engine, UserService *service.UserService, Auth *utils.Auth) {
	userController := &controller.UserController{UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.PUT("/user", userController.UpdateUser)
	}
}

// GetProfileRoutes 获取用户资料
func GetProfileRoutes(router *gin.Engine, ProfileService *service.ProfileService, UserService *service.UserService, Auth *utils.Auth) {
	profileController := &controller.ProfileController{ProfileService: ProfileService, UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.GET("/profiles/:username", profileController.GetProfile)
	}
}

// FollowUserRoutes 关注用户
func FollowUserRoutes(router *gin.Engine, ProfileService *service.ProfileService, UserService *service.UserService, Auth *utils.Auth) {
	profileController := &controller.ProfileController{ProfileService: ProfileService, UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.POST("/profiles/:username/follow", profileController.FollowUser)
	}
}

// UnfollowUserRoutes 取消关注用户
func UnfollowUserRoutes(router *gin.Engine, ProfileService *service.ProfileService, UserService *service.UserService, Auth *utils.Auth) {
	profileController := &controller.ProfileController{ProfileService: ProfileService, UserService: UserService, Auth: Auth}
	api := router.Group("/api")
	{
		api.DELETE("/profiles/:username/follow", profileController.UnfollowUser)
	}
}

// ListArticlesRoutes 文章列表
func ListArticlesRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.GET("/articles", articleController.ListArticles)
	}
}

// FeedArticlesRoutes 关注文章列表
func FeedArticlesRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.GET("/articles/feed", articleController.FeedArticles)
	}
}

// GetArticleRoutes 获取文章
func GetArticleRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.GET("/articles/:slug", articleController.GetArticle)
	}
}

// CreateArticleRoutes 创建文章
func CreateArticleRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.POST("/articles", articleController.CreateArticle)
	}
}

// UpdateArticleRoutes 更新文章
func UpdateArticleRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.PUT("/articles/:slug", articleController.UpdateArticle)
	}
}

// DeleteArticleRoutes 删除文章
func DeleteArticleRoutes(router *gin.Engine, ArticleService *service.ArticleService, Auth *utils.Auth) {
	articleController := &controller.ArticleController{ArticleService: ArticleService, Auth: Auth}
	api := router.Group("/api")
	{
		api.DELETE("/articles/:slug", articleController.DeleteArticle)
	}
}
