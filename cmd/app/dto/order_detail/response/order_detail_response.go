package response

type OrderDetailResponse struct {
	OrderDetailID int32   `json:"orderDetailId"`
	VariantID     int32   `json:"variantId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
	VariantImage  string  `json:"variantImage"`
}
