package models

type OrderDetail struct {
	OrderDetailId int32 `json:"orderDetailId"`
	OrderId       int32 `json:"orderId"`
	ProductId     int32 `json:"productId"`
	Quantity      int32 `json:"quantity"`
	Price         int32 `json:"price"`
}
