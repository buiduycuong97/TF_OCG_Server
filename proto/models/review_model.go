package models

import "time"

type Review struct {
	ReviewID  int32     `gorm:"primaryKey;autoIncrement" json:"reviewID"`
	UserID    int32     `json:"userID"`
	VariantID int32     `json:"variantID"`
	Rating    int32     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
