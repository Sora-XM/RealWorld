package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username  string `gorm:"size:255;unique;not null" json:"username"`
	Bio       string `gorm:"text" json:"bio"`
	Image     string `gorm:"size:255" json:"image"`
	Following bool   `gorm:"-" json:"following"`
}
