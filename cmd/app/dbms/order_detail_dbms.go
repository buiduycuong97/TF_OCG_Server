package dbms

import (
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

type OrderResult struct {
	OrderDetailID int32   `json:"orderDetailId"`
	OrderID       int32   `json:"orderId"`
	Quantity      int32   `json:"quantity"`
	Price         float64 `json:"price"`
	IsReview      bool    `json:"isReview"`

	VariantID    int32  `gorm:"primaryKey;autoIncrement" json:"variantId"`
	ProductID    int32  `json:"productId"`
	Title        string `json:"title"`
	VariantPrice int32  `json:"variantDetail.price"`
	ComparePrice int32  `json:"comparePrice"`
	CountInStock int32  `json:"countInStock"`
	Image        string `json:"image"`
	OptionValue1 int32  `json:"optionValue1"`
	OptionValue2 int32  `json:"optionValue2"`
}

func CreateOrderDetail(tx *gorm.DB, orderDetail *models.OrderDetail) error {
	orderDetail.IsReview = false

	if err := tx.Create(orderDetail).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOrderDetailIsReview(orderDetailID int32, isReview bool) error {
	var orderDetail models.OrderDetail
	result := database.Db.Model(&orderDetail).Where("order_detail_id = ?", orderDetailID).First(&orderDetail)
	if result.Error != nil {
		return result.Error
	}

	orderDetail.IsReview = isReview

	result = database.Db.Save(&orderDetail)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteOrderDetailByVariantID(tx *gorm.DB, variantID int32) error {
	if err := tx.Where("variant_id = ?", variantID).Delete(&models.OrderDetail{}).Error; err != nil {
		return err
	}

	return nil
}

func DeleteOrderDetailByOrderId(tx *gorm.DB, orderID int32) error {
	if err := tx.Where("order_id = ?", orderID).Delete(&models.OrderDetail{}).Error; err != nil {
		return err
	}

	return nil
}
