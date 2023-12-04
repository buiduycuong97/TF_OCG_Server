package dbms

import (
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

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
