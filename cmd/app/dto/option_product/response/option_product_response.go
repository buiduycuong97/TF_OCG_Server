package response

import "tf_ocg/cmd/app/dto/option_value/response"

type OptionProductResponse struct {
	OptionProductID int32                          `json:"optionProductId"`
	ProductID       int32                          `json:"productId"`
	OptionType      string                         `json:"optionType"`
	OptionValues    []response.OptionValueResponse `json:"optionValues"`
}
