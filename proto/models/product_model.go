package models

import "time"

type Product struct {
	ProductID   int32     `json:"productId"`
	Handle      string    `json:"handle"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	CategoryID  int       `json:"categoryID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
