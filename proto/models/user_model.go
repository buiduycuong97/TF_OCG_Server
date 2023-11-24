package models

import (
	"time"
)

type User struct {
	UserID       int32     `gorm:"primaryKey;autoIncrement" json:"userID"`
	UserName     string    `json:"userName"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	UserType     string    `gorm:"default:app" json:"userType"`
	Role         string    `gorm:"default:user" json:"role"`
	RefreshToken string    `json:"refreshToken"`
	ResetToken   string    `json:"resetToken"`
	PhoneNumber  string    `json:"phoneNumber"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
