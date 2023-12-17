package response

import "time"

type ProductInfo struct {
	ProductID         int32     `json:"productId"`
	Handle            string    `json:"handle"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	CategoryID        int       `json:"categoryID"`
	QuantityRemaining int32     `json:"quantityRemaining"`
	Image             string    `json:"image"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
