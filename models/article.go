package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type TagList []string

// Value 将 TagList 转换为 JSON 字节切片
func (t TagList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan 将 JSON 字节切片转换为TagList
func (t *TagList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type Article struct {
	gorm.Model
	Slug           string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `json:"description"`
	Body           string    `gorm:"not null" json:"body"`
	TagList        TagList   `gorm:"type:json" json:"tagList"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	AuthorID       uint      `gorm:"not null" json:"-"`
	Author         UserModel `gorm:"foreignKey:AuthorID" json:"author"`
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

type UpdateArticleRequest struct {
	Article struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Body        *string `json:"body"`
	} `json:"article"`
}
