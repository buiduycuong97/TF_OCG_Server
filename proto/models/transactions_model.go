package models

import "time"

type Transaction struct {
	TransactionID int32     `gorm:"primaryKey;autoIncrement" json:"transactionID"`
	OrderID       int32     `json:"orderID"`
	PaypalOrderID string    `json:"paypalOrderID"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}
