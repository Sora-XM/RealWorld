package service

import (
	"errors"
	"github.com/gosimple/slug"
	"goDemo/models"
	"gorm.io/gorm"
	"time"
)

type ArticleService struct {
	DB *gorm.DB
}

type ListArticlesParams struct {
	Tag       string
	Author    string
	Favorited string
	Limit     int
	Offset    int
}

type FeedArticlesParams struct {
	Limit  int
	Offset int
}

func (s *ArticleService) ListArticles(userID uint, params ListArticlesParams) ([]models.Article, int64, error) {
	query := s.DB.Model(&models.Article{}).Preload("Author").Order("created_at DESC")
	if params.Tag != "" {
		query = query.Where("tag_list @> ?", []string{params.Tag})
	}
	if params.Author != "" {
		query = query.Joins("JOIN authors ON articles.author_id = authors.id").Where("authors.username = ?", params.Author)
	}
	if params.Favorited != "" {
		if userID == 0 {
			return nil, 0, errors.New("未登录用户不能查询收藏文章")
		}
		isFavorited := params.Favorited == "true"
		if isFavorited {
			// 查询被指定用户收藏的文章
			query = query.Joins("JOIN favorites ON articles.id = favorites.article_id").
				Where("favorites.user_id = ?", userID)
		} else {
			// 查询未被指定用户收藏的文章
			query = query.Where("NOT EXISTS (SELECT 1 FROM favorites WHERE favorites.article_id = articles.id AND favorites.user_id = ?)", userID)
		}
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if params.Limit == 0 {
		params.Limit = 20
	}
	query = query.Limit(params.Limit)
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}
	var articles []models.Article
	err = query.Find(&articles).Error
	return articles, total, err
}

func (s *ArticleService) FeedArticles(userID uint, params FeedArticlesParams) ([]models.Article, int64, error) {

	var articles []models.Article
	var count int64

	subQuery := s.DB.Table("follows").
		Select("followed").
		Where("follower = ?", userID)

	// 查询这些ID的文章，并预加载作者信息
	query := s.DB.Model(&models.Article{}).
		Preload("Author").
		Where("author_id IN (?)", subQuery).
		Where("deleted_at IS NULL")

	// 分页查询文章列表
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Limit(params.Limit).Offset(params.Offset).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, count, err
}

func (s *ArticleService) GetArticle(slug string) (*models.Article, error) {
	var article models.Article
	err := s.DB.Preload("Author").Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(userID uint, req models.CreateArticleRequest) (*models.Article, error) {
	slug := GenerateSlug(req.Article.Title)
	article := models.Article{
		Slug:        slug,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
		AuthorID:    userID,
	}
	err := s.DB.Create(&article).Error
	if err != nil {
		return nil, err
	}
	//添加作者
	err = s.DB.Preload("Author").First(&article, "id=?", article.ID).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(userID uint, slug string, req models.UpdateArticleRequest) (*models.Article, error) {
	var article models.Article
	//校验权限
	err := s.DB.Where("slug = ? AND author_id=?", slug, userID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章未找到或无权更新")
		}
		return nil, err
	}
	//更新文章字段
	if req.Article.Title != nil {
		article.Title = *req.Article.Title
		article.Slug = GenerateSlug(*req.Article.Title)
	}
	if req.Article.Description != nil {
		article.Description = *req.Article.Description
	}
	if req.Article.Body != nil {
		article.Body = *req.Article.Body
	}
	err = s.DB.Save(&article).Error
	if err != nil {
		return nil, err
	}
	err = s.DB.Preload("Author").First(&article, "id=?", article.ID).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(userID uint, slug string) error {
	var article models.Article
	err := s.DB.Where("slug = ? AND author_id=?", slug, userID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章未找到或无权删除")
		}
		return err
	}
	err = s.DB.Delete(&article).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateComment 创建评论
func (s *ArticleService) CreateComment(userID uint, slug string, req models.CreateCommentRequest) (*models.Comment, error) {
	var article models.Article
	err := s.DB.Where("slug =?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章未找到")
		}
		return nil, err
	}

	comment := models.Comment{
		Body:      req.Comment.Body,
		AuthorID:  userID,
		ArticleID: article.ID,
	}

	err = s.DB.Create(&comment).Error
	if err != nil {
		return nil, err
	}

	err = s.DB.Preload("Author").First(&comment, "id=?", comment.ID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsBySlug 获取文章的评论列表
func (s *ArticleService) GetCommentsBySlug(slug string, userID uint) ([]models.CommentResponse, error) {
	var article models.Article
	err := s.DB.Where("slug =?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章没找到哦")
		}
		return nil, err
	}
	var comments []models.Comment
	err = s.DB.Preload("Author").Where("article_id = ?", article.ID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	var commentResponses []models.CommentResponse
	for _, comment := range comments {
		isFollowing, err := s.IsFollowing(userID, comment.AuthorID)
		if err != nil {
			return nil, err
		}
		commentResponse := models.CommentResponse{
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
		commentResponses = append(commentResponses, commentResponse)
	}
	return commentResponses, nil
}

// DeleteComment 删除评论
func (s *ArticleService) DeleteComment(userID uint, slug string, commentID uint) error {
	//定位文章评论
	var article models.Article
	err := s.DB.Where("slug =?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章没找到哦")
		}
		return err
	}
	var Comment models.Comment
	err = s.DB.Where("id =? AND article_id =?", commentID, article.ID).First(&Comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论没找到或无权删除")
		}
		return err
	}
	err = s.DB.Delete(&Comment).Error
	if err != nil {
		return err
	}
	return nil
}

// FavoriteArticle 添加文章收藏
func (s *ArticleService) FavoriteArticle(userID uint, slug string) (*models.Article, error) {
	var article models.Article
	err := s.DB.Where("slug =?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章没找到哦")
		}
		return nil, err
	}
	var favorite models.Favorite
	err = s.DB.Where("user_id =? AND article_id =?", userID, article.ID).First(&favorite).Error
	if err == nil {
		// 已收藏，预加载作者信息并返回文章
		err = s.DB.Preload("Author").Where("id = ?", article.ID).First(&article).Error
		if err != nil {
			return nil, err
		}
		return &article, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	favorite = models.Favorite{
		UserID:    userID,
		ArticleID: article.ID,
	}
	err = s.DB.Create(&favorite).Error
	if err != nil {
		return nil, err
	}
	// 更新文章的收藏数
	err = s.DB.Model(&article).Update("favorites_count", gorm.Expr("favorites_count + ?", 1)).Error
	if err != nil {
		return nil, err
	}
	// 预加载作者信息并返回文章
	err = s.DB.Preload("Author").Where("id = ?", article.ID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// UnfavoriteArticle 取消文章收藏
func (s *ArticleService) UnfavoriteArticle(userID uint, slug string) (*models.Article, error) {
	var article models.Article
	err := s.DB.Where("slug =?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章没找到哦")
		}
		return &article, err
	}
	var favorite models.Favorite
	err = s.DB.Where("user_id=? AND article_id =?", userID, article.ID).First(&favorite).Error
	// 未收藏，直接返回
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = s.DB.Preload("Author").Where("id = ?", article.ID).First(&article).Error
		if err != nil {
			return nil, err
		}
		return &article, nil
	} else if err != nil {
		return nil, err
	}
	// 取消收藏并修改收藏数
	err = s.DB.Delete(&favorite).Error
	if err != nil {
		return nil, err
	}
	err = s.DB.Model(&article).Update("favorites_count",
		gorm.Expr("CASE WHEN  favorites_count>0 THEN favorites_count -1 ELSE 0 END")).Error
	if err != nil {
		return nil, err
	}
	err = s.DB.Preload("Author").Where("id = ?", article.ID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (s *ArticleService) IsFollowing(followerID, followedID uint) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Follow{}).
		Where("follower =? AND followed =?", followerID, followedID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GenerateSlug 根据标题生成文章的slug
// 标题："Hello, World!"
// Slug："hello-world"
func GenerateSlug(title string) string {
	//slug := strings.ToLower(title)
	//for _, char := range slug {
	//	if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
	//		slug = strings.ReplaceAll(slug, string(char), "-")
	//	}
	//}
	////移除多余和首尾连字符
	//for strings.Contains(slug, "--") {
	//	slug = strings.ReplaceAll(slug, "--", "-")
	//}
	//slug = strings.Trim(slug, "-")
	//return slug

	// 使用 slug 库生成支持中文的 slug
	return slug.Make(title)
}
