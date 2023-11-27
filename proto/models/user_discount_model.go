package models

type UserDiscount struct {
	UserDiscountID int32 `gorm:"primaryKey;autoIncrement" json:"userDiscountID"`
	UserID         int32 `json:"userID"`
	DiscountID     int32 `json:"discountID"`
}
