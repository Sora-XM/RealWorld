package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username string `gorm:"size:255;unique;not null" json:"username"`
	Email    string `gorm:"size:255;unique;not null" json:"email"`
	Password string `gorm:"text;not null;column:password" json:"-"`
	Bio      string `gorm:"text" json:"bio"`
	Image    string `gorm:"size:255" json:"image"`
}

type RegisterRequest struct {
	User UserModel `json:"user"`
}

type UserLoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

type UserResponse struct {
	User struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}
