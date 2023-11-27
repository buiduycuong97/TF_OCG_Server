package dbms

import (
	"errors"
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CheckUserUsedDiscount(userID int32, discountID int32) (bool, error) {
	var userDiscount models.UserDiscount
	err := database.Db.Where("user_id = ? AND discount_id = ?", userID, discountID).First(&userDiscount).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return errors.Is(err, gorm.ErrRecordNotFound), nil
}

func MarkDiscountAsUsed(userID int32, discountID int32) error {
	userDiscount := models.UserDiscount{
		UserID:     userID,
		DiscountID: discountID,
	}
	return database.Db.Create(&userDiscount).Error
}
