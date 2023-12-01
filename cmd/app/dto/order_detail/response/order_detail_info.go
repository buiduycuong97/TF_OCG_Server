package response

import (
	"tf_ocg/cmd/app/dto/product_dto/response"
	responseuser "tf_ocg/cmd/app/dto/user_dto/response"
	responsevariant "tf_ocg/cmd/app/dto/variant_dto/response"
)

type OrderDetailInfo struct {
	OrderDetailID int32                           `json:"orderDetailId"`
	OrderID       int32                           `json:"orderId"`
	Quantity      int32                           `json:"quantity"`
	Price         float64                         `json:"price"`
	UserInfo      responseuser.UserInfo           `json:"user"`
	ProductInfo   response.ProductInfo            `json:"product"`
	VariantInfo   responsevariant.VariantResponse `json:"variant"`
}
