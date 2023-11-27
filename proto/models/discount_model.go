package models

import "time"

type Discount struct {
	DiscountID        int32     `gorm:"primaryKey;autoIncrement" json:"discountID"`
	DiscountCode      string    `json:"discountCode"`
	DiscountType      string    `json:"discountType"`
	Value             float64   `json:"value"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	AvailableQuantity int32     `json:"availableQuantity"`
}
