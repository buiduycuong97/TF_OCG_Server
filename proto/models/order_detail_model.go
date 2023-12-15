package models

type OrderDetail struct {
	OrderDetailID int32   `gorm:"primaryKey;autoIncrement" json:"orderDetailId"`
	OrderID       int32   `gorm:"foreignKey:OrderID" json:"orderId"`
	VariantID     int32   `gorm:"foreignKey:VariantID" json:"variantId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
	IsReview      bool    `json:"isReview"`
}
