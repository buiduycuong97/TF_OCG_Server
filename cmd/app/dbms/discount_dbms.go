package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateDiscount(discount *models.Discount) (*models.Discount, error) {
	existingDiscount := &models.Discount{}
	database.Db.Raw("SELECT * FROM discounts WHERE discount_code = ?", discount.DiscountCode).Scan(existingDiscount)
	if existingDiscount.DiscountID != 0 {
		return nil, errors.New("Discount code already exists")
	}

	now := time.Now()
	discount.StartDate = now
	discount.EndDate = now.AddDate(0, 1, 0)
	err := database.Db.Create(discount).Error
	if err != nil {
		return nil, err
	}
	return discount, nil
}

func GetDiscountByID(discount *models.Discount, id int32) (err error) {
	err = database.Db.Where("discount_id = ?", id).Find(discount).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllDiscounts() ([]*models.Discount, error) {
	discounts := []*models.Discount{}
	err := database.Db.Find(&discounts).Error
	if err != nil {
		return nil, err
	}
	return discounts, nil
}

func UpdateDiscount(updatedDiscount *models.Discount, id int32) error {
	database.Db.Model(updatedDiscount).Where("discount_id = ?", id).Updates(updatedDiscount)
	return nil
}

func DeleteDiscount(discount *models.Discount, id int32) error {
	database.Db.Where("discount_id = ?", id).Delete(discount)
	return nil
}

func SaveUserDiscount(userID, discountID int32) error {
	userDiscount := &models.UserDiscount{
		UserID:     userID,
		DiscountID: discountID,
	}

	err := database.Db.Create(userDiscount).Error
	if err != nil {
		return err
	}

	return nil
}
