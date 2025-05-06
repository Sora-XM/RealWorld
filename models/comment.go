package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	Body      string    `gorm:"not null" json:"body"`
	AuthorID  uint      `gorm:"not null" json:"-"`
	Author    UserModel `gorm:"foreignKey:AuthorID" json:"author"`
	ArticleID uint      `gorm:"not null" json:"-"`
}

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment" binding:"required"`
}

type CommentResponse struct {
	Comment struct {
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
	} `json:"comment"`
}

type CommentsResponse struct {
	Comments []struct {
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
	} `json:"comments"`
}
