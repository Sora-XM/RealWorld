package service

import (
	"errors"
	"goDemo/models"
	"gorm.io/gorm"
	"strings"
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

func (s *ArticleService) ListArticles(params ListArticlesParams) ([]models.Article, int64, error) {
	query := s.DB.Model(&models.Article{}).Preload("Author").Order("created_at DESC")

	if params.Tag != "" {
		query = query.Where("tag_list @> ?", []string{params.Tag})
	}
	if params.Author != "" {
		query = query.Joins("JOIN authors ON articles.author_id = authors.id").Where("authors.username = ?", params.Author)
	}
	if params.Favorited != "" {
		// 暂时简化处理
		query = query.Where("favorited = true")
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

	//先获取用户关注的作者ID
	subQuery := s.DB.Table("follows").
		Select("followed").
		Where("follower = ?", userID)

	//查询这些ID的文章
	query := s.DB.Model(&models.Article{}).
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
	err := s.DB.Where("slug = ?", slug).First(&article).Error
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

// GenerateSlug 根据标题生成文章的slug
// 标题："Hello, World!"
// Slug："hello-world"
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			slug = strings.ReplaceAll(slug, string(char), "-")
		}
	}
	//移除多余和首尾连字符
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	return slug
}
