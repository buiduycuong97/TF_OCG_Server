package response

type Product struct {
	Handle            string  `json:"handle"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Price             float64 `json:"price"`
	CategoryID        int     `json:"categoryID"`
	QuantityRemaining int32   `json:"quantity_remaining"`
}
