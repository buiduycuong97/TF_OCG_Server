package models

import "time"

type OrderStatus string

const (
	Unpaid           OrderStatus = "unpaid"
	Paid             OrderStatus = "paid"
	BeingTransported OrderStatus = "being transported"
	PreparingOrders  OrderStatus = "preparing orders"
	Delivered        OrderStatus = "delivered"
)

type Order struct {
	OrderID         int32       `json:"orderId"`
	UserId          int32       `json:"userId"`
	OrderDate       time.Time   `json:"orderDate"`
	ShippingAddress string      `json:"shippingAddress"`
	Status          OrderStatus `gorm:"default:unpaid" json:"status"`
}
