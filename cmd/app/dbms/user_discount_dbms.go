package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

//func CheckUserUsedDiscount(userID int32, discountID int32) (bool, error) {
//	var userDiscount models.UserDiscount
//	err := database.Db.Where("user_id = ? AND discount_id = ?", userID, discountID).First(&userDiscount).Error
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return false, err
//	}
//	return errors.Is(err, gorm.ErrRecordNotFound), nil
//}
//
//func MarkDiscountAsUsed(userID int32, discountID int32) error {
//	userDiscount := models.UserDiscount{
//		UserID:     userID,
//		DiscountID: discountID,
//	}
//	return database.Db.Create(&userDiscount).Error
//}

func CreateUserDiscount(userDiscount *models.UserDiscount) error {
	if userDiscount.UserID == 0 || userDiscount.DiscountID == 0 {
		return errors.New("Invalid UserID or DiscountID")
	}

	user := &models.User{}
	if err := database.Db.First(user, userDiscount.UserID).Error; err != nil {
		return errors.New("User not found")
	}

	discount := &models.Discount{}
	if err := database.Db.First(discount, userDiscount.DiscountID).Error; err != nil {
		return errors.New("Discount not found")
	}

	if err := database.Db.Create(userDiscount).Error; err != nil {
		return err
	}

	return nil
}

func GetUserDiscountByDiscountCodeAndUserID(userDiscount *models.UserDiscount, discountCode string, userID int) error {
	err := database.Db.
		Joins("JOIN discounts ON user_discounts.discount_id = discounts.discount_id").
		Where("discounts.discount_code = ? AND user_discounts.user_id = ?", discountCode, userID).
		First(userDiscount).
		Error

	return err
}
