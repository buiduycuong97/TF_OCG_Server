package models

type OrderDetail struct {
	OrderDetailID int32   `gorm:"primaryKey;autoIncrement" json:"orderDetailId"`
	OrderID       int32   `json:"orderId"`
	ProductID     int32   `json:"productId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
}
