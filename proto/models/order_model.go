package models

import "time"

type OrderStatus string

const (
	Pending              OrderStatus = "pending"
	OrderBeingDelivered  OrderStatus = "order being delivered"
	CompleteTheOrder     OrderStatus = "complete the order"
	RequestToCancelOrder OrderStatus = "request to cancel order"
	Cancelled            OrderStatus = "cancelled"
)

type Order struct {
	OrderID         int32       `gorm:"primaryKey;autoIncrement" json:"orderId"`
	UserID          int32       `json:"userId"`
	OrderDate       time.Time   `json:"orderDate"`
	ShippingAddress string      `json:"shippingAddress"`
	Status          OrderStatus `gorm:"default:unpaid" json:"status"`
}
