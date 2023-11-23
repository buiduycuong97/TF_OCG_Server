package dbms

import (
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func ClearCart(tx *gorm.DB, userID int32) error {
	cartItems := []*models.Cart{}
	if err := tx.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return err
	}

	for _, cartItem := range cartItems {
		if err := tx.Delete(cartItem).Error; err != nil {
			return err
		}
	}

	return nil
}

func CreateOrder(tx *gorm.DB, order *models.Order) (*models.Order, error) {
	now := time.Now()
	order.OrderDate = now
	order.Status = models.Pending

	err := tx.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func UpdateOrderStatus(orderID int32, status string) error {
	err := database.Db.Model(&models.Order{}).Where("order_id = ?", orderID).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOrderDetailsByOrderID(orderID int32) ([]models.OrderDetail, error) {
	var orderDetails []models.OrderDetail

	result := database.Db.Where("order_id = ?", orderID).Find(&orderDetails)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderDetails, nil
}

func GetOrderStatus(orderID int32) (string, error) {
	var order models.Order
	if err := database.Db.First(&order, orderID).Error; err != nil {
		return "", err
	}

	return string(order.Status), nil
}

func GetOrdersByStatus(status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order

	result := database.Db.Where("status = ?", status).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}
