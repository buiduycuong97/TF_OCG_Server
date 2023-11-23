package dbms

import (
	"gorm.io/gorm"
	"tf_ocg/proto/models"
)

func CreateOrderDetail(tx *gorm.DB, orderDetail *models.OrderDetail) error {
	if err := tx.Create(orderDetail).Error; err != nil {
		return err
	}
	return nil
}
