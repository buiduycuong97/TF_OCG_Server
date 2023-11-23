package response

type OrderDetailResponse struct {
	OrderDetailID int32   `json:"orderDetailId"`
	ProductID     int32   `json:"productId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
}
