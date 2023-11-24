package response

import "tf_ocg/proto/models"

type CartResponse struct {
	CartID        int32          `json:"cartId"`
	UserID        int32          `json:"userId"`
	ProductID     int32          `json:"productId"`
	Quantity      int32          `json:"quantity"`
	TotalPrice    float64        `json:"totalPrice"`
	ProductDetail models.Product `json:"productDetail"`
}
