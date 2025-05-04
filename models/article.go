package models

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Slug           string   `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Title          string   `gorm:"not null" json:"title"`
	Description    string   `json:"description"`
	Body           string   `gorm:"not null" json:"body"`
	TagList        []string `gorm:"type:json" json:"tagList"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	AuthorID       uint     `gorm:"not null" json:"-"`
	Author         Author   `gorm:"foreignKey:AuthorID" json:"author"`
}

type Author struct {
	gorm.Model
	Username  string `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ArticleListResponse struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

type ArticleResponse struct {
	Article Article `json:"article"`
}

type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Body        string   `json:"body" binding:"required"`
		TagList     []string `json:"tagList"`
	} `json:"article" binding:"required"`
}
