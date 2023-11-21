package request

type CreateProductReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}
