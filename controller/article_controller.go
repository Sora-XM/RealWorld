package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goDemo/models"
	"goDemo/service"
	"goDemo/utils"
	"net/http"
)

type ArticleController struct {
	ArticleService *service.ArticleService
	Auth           *utils.Auth
}

// ListArticles 文章列表
// @Summary 文章列表
// @Description 获取文章列表
// @Tags articles
// @Accept json
// @Produce json
// @Param tag query string false "标签"
// @Param author query string false "作者"
// @Param favorited query string false "是否收藏"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} models.ArticleListResponse
// @Router /api/articles [get]
func (c *ArticleController) ListArticles(ctx *gin.Context) {
	param := service.ListArticlesParams{
		Tag:       ctx.Query("tag"),
		Author:    ctx.Query("author"),
		Favorited: ctx.Query("favorited"),
		Limit:     getIntQuery(ctx, "limit", 20),
		Offset:    getIntQuery(ctx, "offset", 0),
	}
	articles, count, err := c.ArticleService.ListArticles(param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	ctx.JSON(http.StatusOK, models.ArticleListResponse{
		Articles:      articles,
		ArticlesCount: int(count),
	})
}

// FeedArticles 关注文章列表
// @Summary 关注文章列表
// @Description 获取关注文章列表
// @Tags articles
// @Accept json
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} models.ArticleListResponse
// @Router /api/articles/feed [get]
func (c *ArticleController) FeedArticles(ctx *gin.Context) {
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}
	param := service.FeedArticlesParams{
		Limit:  getIntQuery(ctx, "limit", 20),
		Offset: getIntQuery(ctx, "offset", 0),
	}
	articles, count, err := c.ArticleService.FeedArticles(userID, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	ctx.JSON(http.StatusOK, models.ArticleListResponse{
		Articles:      articles,
		ArticlesCount: int(count),
	})
}

// GetArticle 获取文章
// @Summary 获取文章
// @Description 获取文章详情
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "文章slug"
// @Success 200 {object} models.ArticleResponse
// @Router /api/articles/{slug} [get]
func (c *ArticleController) GetArticle(ctx *gin.Context) {

	slug := ctx.Param("slug")
	article, err := c.ArticleService.GetArticle(slug)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	if article == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"文章没找到哦"}}})
		return
	}
	ctx.JSON(http.StatusOK, models.ArticleResponse{Article: *article})
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags articles
// @Accept json
// @Produce json
// @Param article body models.CreateArticleRequest true "文章信息"
// @Success 201 {object} models.ArticleResponse
// @Router /api/articles [post]
func (c *ArticleController) CreateArticle(ctx *gin.Context) {

	//校验token
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}

	//解析请求内容
	var request models.CreateArticleRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}

	//调用服务层创建文章
	article, err := c.ArticleService.CreateArticle(userID, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}

	//返回创建成功的文章信息
	ctx.JSON(http.StatusOK, models.ArticleResponse{Article: *article})
}

// 获取查询参数并设置默认值
// 在没有获取到查询参数时，使用默认值
func getIntQuery(ctx *gin.Context, key string, defaultValue int) int {
	value, err := ctx.GetQuery(key)
	if err {
		var num int
		fmt.Sscanf(value, "%d", &num)
		return num
	}
	return defaultValue
}
