package response

import (
	"tf_ocg/cmd/app/dto/option_product/response"
	"tf_ocg/proto/models"
)

type ProductWithOptionResponse struct {
	Product        models.Product                   `json:"product"`
	OptionProducts []response.OptionProductResponse `json:"optionProducts"`
	Variants       []models.Variant                 `json:"variants"`
}
