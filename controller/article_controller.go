package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goDemo/models"
	"goDemo/service"
	"goDemo/utils"
	"gorm.io/gorm"
	"net/http"
	"time"
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
// @Security BearerAuth
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
// @Security BearerAuth
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
// @Security BearerAuth
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

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "文章slug"
// @Param article body models.UpdateArticleRequest true "文章信息"
// @Success 200 {object} models.ArticleResponse
// @Router /api/articles/{slug} [put]
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	//校验token
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}
	//解析请求内容并调用服务层更新文章
	slug := ctx.Param("slug")
	var request models.UpdateArticleRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	article, err := c.ArticleService.UpdateArticle(userID, slug, request)
	//处理错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"文章没找到或无权查看"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	//返回更新成功的文章信息
	ctx.JSON(http.StatusOK, models.ArticleResponse{Article: *article})
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 删除指定文章
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "文章slug"
// @Success 204 {object} models.ArticleResponse
// @Router /api/articles/{slug} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	//校验token
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}
	//调用服务层删除文章
	slug := ctx.Param("slug")
	err = c.ArticleService.DeleteArticle(userID, slug)
	//处理错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"文章未找到或无权限删除"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}

// AddComment 向文章添加评论
// @Summary 添加评论
// @Description 向文章添加评论
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "文章slug"
// @Param comment body models.CreateCommentRequest true "评论信息"
// @Success 201 {object} models.CommentResponse
// @Router /api/articles/{slug}/comments [post]
func (c *ArticleController) AddComment(ctx *gin.Context) {
	//校验token
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}
	slug := ctx.Param("slug")
	var request models.CreateCommentRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
	}
	//创建评论并返回信息
	comment, err := c.ArticleService.CreateComment(userID, slug, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"文章未找到"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	//是否已关注该评论作者
	isFollowing, err := c.ArticleService.IsFollowing(userID, comment.AuthorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	//构造响应
	response := models.CommentResponse{
		Comment: struct {
			ID        uint      `json:"id"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
			Body      string    `json:"body"`
			Author    struct {
				Username  string `json:"username"`
				Bio       string `json:"bio"`
				Image     string `json:"image"`
				Following bool   `json:"following"`
			} `json:"author"`
		}{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Body:      comment.Body,
			Author: struct {
				Username  string `json:"username"`
				Bio       string `json:"bio"`
				Image     string `json:"image"`
				Following bool   `json:"following"`
			}{
				Username:  comment.Author.Username,
				Bio:       comment.Author.Bio,
				Image:     comment.Author.Image,
				Following: isFollowing,
			},
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// GetComments 获取文章的评论列表
// @Summary 获取文章的评论列表
// @Description 获取指定文章的评论列表
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "文章slug"
// @Success 200 {object} models.CommentsResponse
// @Router /api/articles/{slug}/comments [get]
func (c *ArticleController) GetComments(ctx *gin.Context) {
	// 校验token
	userID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权访问"}}})
		return
	}
	slug := ctx.Param("slug")
	commentResponses, err := c.ArticleService.GetCommentsBySlug(slug, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	// 构造响应
	var commentsResponse models.CommentsResponse
	for _, commentResponse := range commentResponses {
		commentsResponse.Comments = append(commentsResponse.Comments, commentResponse.Comment)
	}
	ctx.JSON(http.StatusOK, commentsResponse)
}

// DeleteComment 删除文章评论
// @Summary 删除文章评论
// @Description 删除指定文章的评论
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "文章slug"
// @Param id path int true "评论ID"
// @Success 204 {object} models.CommentResponse
// @Router /api/articles/{slug}/comments/{id} [delete]
func (c *ArticleController) DeleteComment(ctx *gin.Context) {

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
