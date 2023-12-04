package response

type OptionValueResponse struct {
	OptionValueID   int32  `json:"optionValueId"`
	OptionProductID int32  `json:"optionProductId"`
	Value           string `json:"value"`
}
