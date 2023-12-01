package response

import "tf_ocg/proto/models"

type VariantDetail struct {
	VariantID    int32              `json:"variantId"`
	ProductID    int32              `json:"productId"`
	Title        string             `json:"title"`
	Price        int32              `json:"price"`
	ComparePrice int32              `json:"comparePrice"`
	CountInStock int32              `json:"countInStock"`
	Image        string             `json:"image"`
	OptionValue1 models.OptionValue `json:"optionValue1"`
	OptionValue2 models.OptionValue `json:"optionValue2"`
}
