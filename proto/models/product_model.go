package models

import "time"

type Product struct {
	ProductID         int32     `gorm:"primaryKey;autoIncrement" json:"productId"`
	Handle            string    `json:"handle"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	CategoryID        int       `json:"categoryID"`
	QuantityRemaining int32     `json:"quantity_remaining"`
	Image             string    `json:"image"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
