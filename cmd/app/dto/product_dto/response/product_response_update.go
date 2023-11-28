package response

import "time"

type ProductResponseUpdate struct {
	ProductID         int32     `json:"productId"`
	Handle            string    `json:"handle"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	CategoryID        int       `json:"categoryID"`
	QuantityRemaining int32     `json:"quantityRemaining"`
	UpdatedAt         time.Time `json:"updated_at"`
}
