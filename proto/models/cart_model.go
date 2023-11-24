package models

type Cart struct {
	CartID     int32   `gorm:"primaryKey;autoIncrement" json:"cartId"`
	UserID     int32   `json:"userId"`
	ProductID  int32   `json:"productId"`
	Quantity   int32   `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}
