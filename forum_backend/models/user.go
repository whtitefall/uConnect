package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`

	Posts    []Post    `json:"posts"`
	Comments []Comment `json:"comments"`
}
