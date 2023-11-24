package response

import (
	"tf_ocg/cmd/app/dto/order_detail/response"
	"tf_ocg/proto/models"
	"time"
)

type OrderResponse struct {
	OrderID         int32                          `json:"orderId"`
	UserID          int32                          `json:"userId"`
	OrderDate       time.Time                      `json:"orderDate"`
	ShippingAddress string                         `json:"shippingAddress"`
	Status          models.OrderStatus             `json:"status"`
	OrderDetails    []response.OrderDetailResponse `json:"orderDetails"`
}
