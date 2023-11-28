package models

type OptionValue struct {
	OptionValueID   int32  `gorm:"primaryKey;autoIncrement" json:"optionValueId"`
	OptionProductID int32  `json:"optionProductId"`
	Value           string `json:"value"`
}
