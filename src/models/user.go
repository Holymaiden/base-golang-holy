package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"AUTO_INCREMENT" json:"id"`
	Name      string         `json:"name" binding:"required" gorm:"type:varchar(50)"`
	Username  string         `gorm:"unique;type:varchar(50)" json:"username" binding:"required"`
	Email     string         `gorm:"unique;type:varchar(50)" json:"email" binding:"required,email"`
	Password  string         `json:"password" binding:"required" gorm:"type:varchar(255)"`
	Role      string         `json:"role" binding:"required"`
	Avatar    string         `json:"avatar"`
	Token     string         `gorm:"type:varchar(255)" json:"token,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
