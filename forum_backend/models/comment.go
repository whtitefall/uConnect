package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content"`

	UserID uint `json:"user_id"`
	User   User `json:"user"` // 自动 preload

	PostID uint `json:"post_id"`
	Post   Post `json:"-"` // 避免递归嵌套
}
