package response

type Product struct {
	ProductID   int32  `json:"productId"`
	Handle      string `json:"handle"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}
