package response

import (
	"tf_ocg/cmd/app/dto/variant_dto/response"
	"tf_ocg/proto/models"
)

type CartResponse struct {
	CartID        int32                  `json:"cartId"`
	UserID        int32                  `json:"userId"`
	ProductID     int32                  `json:"productId"`
	VariantID     int32                  `json:"variantId"`
	Quantity      int32                  `json:"quantity"`
	TotalPrice    float64                `json:"totalPrice"`
	ProductDetail models.Product         `json:"productDetail"`
	VariantDetail response.VariantDetail `json:"variantDetail"`
}
