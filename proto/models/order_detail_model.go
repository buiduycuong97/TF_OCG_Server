package models

type OrderDetail struct {
	OrderDetailID int32   `gorm:"primaryKey;autoIncrement" json:"orderDetailId"`
	OrderID       int32   `json:"orderId"`
	VariantID     int32   `json:"variantId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
	IsReview      bool    `json:"isReview"`
}
