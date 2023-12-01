package response

type VariantResponse struct {
	VariantID      int32  `gorm:"primaryKey;autoIncrement" json:"variantId"`
	ProductID      int32  `json:"productId"`
	Title          string `json:"title"`
	Price          int32  `json:"price"`
	ComparePrice   int32  `json:"comparePrice"`
	CountInStock   int32  `json:"countInStock"`
	OptionProduct1 int32  `json:"optionProduct1"`
	OptionProduct2 int32  `json:"optionProduct2"`
}
