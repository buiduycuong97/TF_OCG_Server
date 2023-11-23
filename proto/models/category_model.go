package models

type Categories struct {
	CategoryID int32  `gorm:"primaryKey;autoIncrement" json:"categoryId"`
	Handle     string `json:"handle"`
	Name       string `json:"name"`
	Image      string `json:"image"`
}
