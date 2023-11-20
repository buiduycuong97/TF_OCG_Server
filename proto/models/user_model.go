package models

import (
	"time"
)

type User struct {
	UserID       int32     `json:"userID"`
	UserName     string    `json:"userName"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	UserType     string    `gorm:"default:app" json:"userType"`
	Role         string    `gorm:"default:user" json:"role"`
	RefreshToken string    `json:"refreshToken"`
	ResetToken   string    `json:"resetToken"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
