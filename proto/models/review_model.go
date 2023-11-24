package models

type Review struct {
	ReviewID  int32  `json:"reviewID"`
	UserID    int32  `json:"userID"`
	ProductID int32  `json:"productID"`
	Rating    int32  `json:"rating"`
	Comment   string `json:"comment"`
}
