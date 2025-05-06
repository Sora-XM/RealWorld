package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint `gorm:"not null;index" json:"user_id"`    // 用户 ID
	ArticleID uint `gorm:"not null;index" json:"article_id"` // 文章 ID
}
