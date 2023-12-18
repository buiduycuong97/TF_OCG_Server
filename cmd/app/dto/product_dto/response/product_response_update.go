package response

import "time"

type ProductResponseUpdate struct {
	ProductID   int32     `json:"productId"`
	Handle      string    `json:"handle"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  int       `json:"categoryID"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
