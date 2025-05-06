package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goDemo/config"
	_ "goDemo/docs" // 导入生成的文档包
	"goDemo/route"
	"goDemo/service"
	"goDemo/utils"
	"log"
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
	//项目启动方式：
	//1. 启动数据库服务，在config/sql_config处修改数据库连接URL
	//2. 启动项目，访问 http://localhost:8080/swagger/index.html 查看API文档
	//3. 注册账号，登录获取token，在Authorization处填写token，即可访问其他接口
	//Tips：目前登录接口因为自己电脑不知名原因密码校验一直校验失败，
	//因此注释了那段代码，只要输入正确用户名即可成功登录，登录获取token，然后访问其他接口
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	//JWT密钥
	auth := utils.NewAuth("DurRDDtjL2uB_Zyry4f6GHwoBgD5k7oLvC7Fj12E56E=")
	// 初始化服务
	userService := &service.UserService{
		DB:   db,
		Auth: auth,
	}
	profileService := &service.ProfileService{
		DB: db,
	}
	articleService := &service.ArticleService{
		DB: db,
	}
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	// 注册 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	route.SetupRoutes(router, userService, auth)
	route.LoginRoutes(router, userService, auth)
	route.GetCurrentUserRoutes(router, userService, auth)
	route.UpdateUserRoutes(router, userService, auth)
	route.GetProfileRoutes(router, profileService, userService, auth)
	route.FollowUserRoutes(router, profileService, userService, auth)
	route.UnfollowUserRoutes(router, profileService, userService, auth)
	route.ListArticlesRoutes(router, articleService, auth)
	route.FeedArticlesRoutes(router, articleService, auth)
	route.GetArticleRoutes(router, articleService, auth)
	route.CreateArticleRoutes(router, articleService, auth)
	route.UpdateArticleRoutes(router, articleService, auth)
	route.DeleteArticleRoutes(router, articleService, auth)
	route.AddCommentRoutes(router, articleService, auth)
	route.GetCommentsRoutes(router, articleService, auth)
	route.DeleteCommentRoutes(router, articleService, auth)
	route.FavoriteArticleRoutes(router, articleService, auth)
	route.UnfavoriteArticleRoutes(router, articleService, auth)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
