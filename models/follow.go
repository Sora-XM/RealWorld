package models

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	Follower uint `gorm:"index;not null" json:"follower"` // 关注者 ID
	Followed uint `gorm:"index;not null" json:"followed"` // 被关注者 ID
}
