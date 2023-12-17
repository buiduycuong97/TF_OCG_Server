package models

import "time"

type Product struct {
	ProductID   int32     `gorm:"primaryKey;autoIncrement" json:"productId"`
	Handle      string    `json:"handle"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  int       `json:"categoryID"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
