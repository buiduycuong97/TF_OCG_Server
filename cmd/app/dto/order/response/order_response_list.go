package response

import (
	order_detail_response "tf_ocg/cmd/app/dto/order_detail/response"
	"tf_ocg/proto/models"
	"time"
)

type OrderResponseList struct {
	OrderID         int32                                       `json:"orderId"`
	UserID          int32                                       `json:"userId"`
	OrderDate       time.Time                                   `json:"orderDate"`
	ShippingAddress string                                      `json:"shippingAddress"`
	Status          models.OrderStatus                          `json:"status"`
	OrderDetails    []order_detail_response.OrderDetailResponse `json:"orderDetails"`
	TotalPrice      float64                                     `json:"totalPrice"`
}
