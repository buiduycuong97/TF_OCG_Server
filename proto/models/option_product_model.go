package models

type OptionProduct struct {
	OptionProductID int32  `gorm:"primaryKey;autoIncrement" json:"optionProductId"`
	ProductID       int32  `json:"productId"`
	OptionType      string `json:"optionType"`
}
