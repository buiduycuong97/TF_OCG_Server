package response

import (
	"tf_ocg/proto/models"
	"time"
)

type OrderInfo struct {
	OrderID         int32              `json:"orderId"`
	UserID          int32              `json:"userId"`
	OrderDate       time.Time          `json:"orderDate"`
	ShippingAddress string             `json:"shippingAddress"`
	PhoneOrder      string             `json:"phoneOrder"`
	Status          models.OrderStatus `json:"status"`
	ProvinceID      int32              `json:"provinceId"`
	TotalQuantity   int32              `json:"totalQuantity"`
	TotalPrice      float64            `json:"totalPrice"`
	DiscountAmount  float64            `json:"discountAmount"`
	GrandTotal      float64            `json:"grandTotal"`
	OrderDetails    []OrderDetailInfo  `json:"orderDetails"`
	TotalPages      int                `json:"totalPages"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}
