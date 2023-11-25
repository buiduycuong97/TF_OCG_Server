package dbms

import (
	"errors"
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

func UpdateOrderTotalValues(db *gorm.DB, orderID int32, totalQuantity int32, totalPrice float64) error {
	return db.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"TotalQuantity": totalQuantity,
			"TotalPrice":    totalPrice,
		}).Error
}

func ApplyDiscountForOrder(db *gorm.DB, cartItems []models.Cart, discountCode string) (float64, error) {
	discount, err := GetDiscountByCode(db, discountCode)
	if err != nil {
		return 0, err
	}
	if discount.AvailableQuantity <= 0 {
		return 0, errors.New("Discount code is not available")
	}
	var totalDiscount float64
	for _, cartItem := range cartItems {
		productDiscount := calculateProductDiscount(cartItem.TotalPrice, discount)
		cartItem.TotalPrice -= productDiscount
		totalDiscount += productDiscount
	}
	err = decreaseDiscountQuantity(db, discount)
	if err != nil {
		return 0, err
	}
	return totalDiscount, nil
}

func calculateProductDiscount(productPrice float64, discount models.Discount) float64 {
	switch discount.DiscountType {
	case "percentage":
		return productPrice * discount.Value / 100
	case "fixed":
		return discount.Value
	default:
		return 0
	}
}

func decreaseDiscountQuantity(db *gorm.DB, discount models.Discount) error {
	discount.AvailableQuantity--
	return db.Save(&discount).Error
}

func GetDiscountByCode(db *gorm.DB, discountCode string) (models.Discount, error) {
	var discount models.Discount
	err := db.Where("discount_code = ?", discountCode).First(&discount).Error
	if err != nil {
		return discount, err
	}
	return discount, nil
}
