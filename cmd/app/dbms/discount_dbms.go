package dbms

import (
	"errors"
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateDiscount(discount *models.Discount) (*models.Discount, error) {
	existingDiscount := &models.Discount{}
	database.Db.Raw("SELECT * FROM discounts WHERE discount_code = ?", discount.DiscountCode).Scan(existingDiscount)
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

func GetDiscountByDiscountCode(discount *models.Discount, code string) (err error) {
	err = database.Db.Where("discount_code = ?", code).Find(discount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid discount code")
		}
		return err
	}
	return nil
}

func GetAllDiscounts(page int32, pageSize int32, searchText string) ([]*models.Discount, error) {
	discounts := []*models.Discount{}
	offset := (page - 1) * pageSize
	query := database.Db
	query = query.Order("start_date desc")
	query = query.Offset(int(offset)).Limit(int(pageSize))
	if searchText != "" {
		query = query.Where("discount_code LIKE ?", "%"+searchText+"%")
	}

	err := query.Find(&discounts).Error
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
	var userDiscount models.UserDiscount
	if err := database.Db.Where("discount_id = ?", id).Find(&userDiscount).Error; err == nil {
		if err := database.Db.Delete(&userDiscount).Error; err != nil {
			return err
		}
	}

	if err := database.Db.Where("discount_id = ?", id).Delete(discount).Error; err != nil {
		return err
	}

	return nil
}

func DeleteDiscountIfExists(discountID int32) error {
	// Kiểm tra xem có tồn tại mã trong user_discount không
	var userDiscount models.UserDiscount

	// Kiểm tra xem có tồn tại mã trong user_discount không và lấy giá trị UserDiscount nếu tồn tại
	if err := database.Db.Where("discount_id = ?", discountID).First(&userDiscount).Error; err != nil {
		// Nếu không tìm thấy user_discount, chỉ xóa discount
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := database.Db.Exec("DELETE FROM discounts WHERE discount_id = ?", discountID).Error; err != nil {
				return err
			}
			return nil
		}
		// Nếu có lỗi khác, trả về lỗi
		return err
	}

	// Nếu tồn tại trong user_discount, thì xóa user_discount
	if err := database.Db.Delete(&userDiscount).Error; err != nil {
		return err
	}

	// Xóa discount
	if err := database.Db.Exec("DELETE FROM discounts WHERE discount_id = ?", discountID).Error; err != nil {
		return err
	}

	return nil
}

func GetUserDiscountByDiscountID(discountID int32) (*models.UserDiscount, error) {
	var userDiscount models.UserDiscount

	if err := database.Db.Where("discount_id = ?", discountID).First(&userDiscount).Error; err != nil {
		return nil, err
	}

	return &userDiscount, nil
}

func DeleteDiscountAutoGen(discount *models.Discount, id int32) error {
	if err := database.Db.Where("discount_id = ?", id).Delete(discount).Error; err != nil {
		return err
	}

	return nil
}

func GetDiscountByDifferentCode(discount *models.Discount, discountCode string) error {
	err := database.Db.
		Raw(`
			SELECT * FROM discounts
			WHERE discount_code = ?
			AND NOT EXISTS (
				SELECT 1 FROM user_discounts
				WHERE user_discounts.discount_id = discounts.discount_id
			)
		`, discountCode).
		First(discount).
		Error

	return err
}

func GetDiscountByDiscountCodeAndUserID(discount *models.Discount, discountCode string, userID int) error {
	// Kiểm tra giá trị trong user_discounts theo user_id và discount_code
	userDiscountErr := database.Db.
		Joins("JOIN user_discounts ON user_discounts.discount_id = discounts.discount_id").
		Where("discounts.discount_code = ? AND user_discounts.user_id = ?", discountCode, userID).
		First(discount).
		Error

	if userDiscountErr == nil {
		// Giá trị tồn tại trong user_discounts, trả về mà không cần kiểm tra tiếp discounts
		return nil
	}

	// Nếu giá trị không tồn tại trong user_discounts, kiểm tra trong discounts
	discountErr := database.Db.
		Raw(`
		SELECT * FROM discounts
		WHERE discount_code = ?
		AND NOT EXISTS (
			SELECT 1 FROM user_discounts
			WHERE user_discounts.discount_id = discounts.discount_id
		)
	`, discountCode).
		First(discount).
		Error
	// Trả về giá trị discount ngay cả khi có lỗi
	return discountErr
}
