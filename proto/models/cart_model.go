package models

type Cart struct {
	CartID     int32   `gorm:"primaryKey;autoIncrement" json:"cartId"`
	UserID     int32   `json:"userId"`
	VariantID  int32   `json:"variantId"`
	Quantity   int32   `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}
