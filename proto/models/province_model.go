package models

type Province struct {
	ProvinceID   int32   `gorm:"primaryKey;autoIncrement" json:"provinceId"`
	ProvinceName string  `json:"provinceName"`
	ShippingFee  float64 `gorm:"default:0" json:"shippingFee"`
}
