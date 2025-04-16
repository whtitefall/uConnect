package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"` // 外键
	User    User   `json:"user"`    // 自动 preload 用户信息
}
