package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique" json:"username"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"-"`
	Role      string `gorm:"default:user"` // "user" 或 "admin"
	CreatedAt time.Time

	Posts    []Post    `json:"posts"`
	Comments []Comment `json:"comments"`
}
