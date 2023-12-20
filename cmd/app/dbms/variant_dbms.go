package dbms

import (
	"errors"
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateVariant(variant *models.Variant) (*models.Variant, error) {
	existingVariant := models.Variant{}
	database.Db.Raw("SELECT * FROM variants WHERE product_id = ? AND option_value1 = ? AND option_value2 = ?", variant.ProductID, variant.OptionValue1, variant.OptionValue2).Scan(&existingVariant)
	if existingVariant.VariantID > 0 {
		return nil, errors.New("Variant with the same product, optionProduct1, and optionProduct2 already exists")
	}
	err := database.Db.Create(variant).Error
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func GetVariantIdByOption(productID, optionValue1, optionValue2 int32) (int32, error) {
	var variantID int32

	if optionValue2 == 0 {
		result := database.Db.Raw("SELECT variant_id FROM variants WHERE product_id = ? AND option_value1 = ? AND option_value2 = 0", productID, optionValue1).Scan(&variantID)
		if result.Error != nil {
			return 0, result.Error
		}
	} else {
		result := database.Db.Raw("SELECT variant_id FROM variants WHERE product_id = ? AND option_value1 = ? AND option_value2 = ?", productID, optionValue1, optionValue2).Scan(&variantID)
		if result.Error != nil {
			return 0, result.Error
		}
	}

	return variantID, nil
}

func GetVariantById(variant *models.Variant, variantID int32) error {
	err := database.Db.Where("variant_id = ?", variantID).First(variant).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVariantByIdForDelete(tx *gorm.DB, variant *models.Variant, variantID int32) error {
	if err := tx.Where("variant_id = ?", variantID).First(variant).Error; err != nil {
		return err
	}
	return nil
}

func GetVariantsByProductId(productId int32) ([]models.Variant, error) {
	var variants []models.Variant
	err := database.Db.Where("product_id = ?", productId).Find(&variants).Error
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func GetVariantByIdInGetOrder(variantID int32) (models.Variant, error) {
	var variant models.Variant
	err := database.Db.First(&variant, variantID).Error
	return variant, err
}

func UpdateVariant(variant *models.Variant, variantID int32) error {
	return database.Db.Model(&models.Variant{}).Where("variant_id = ?", variantID).Updates(variant).Error
}

func GetVariantIDsByProductID(productID int32) ([]int32, error) {
	var variantIDs []int32

	result := database.Db.Model(&models.Variant{}).Where("product_id = ?", productID).Pluck("variant_id", &variantIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	return variantIDs, nil
}

func UpdateVariantByAdmin(variant *models.Variant, variantID int32) error {
	err := database.Db.Model(&models.Variant{}).Where("variant_id = ?", variantID).Updates(variant).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteVariant(variantID int32) error {
	tx := database.Db.Begin()

	// Xóa variant và các liên kết
	err := GetVariantByIdForDelete(tx, &models.Variant{}, variantID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := DeleteReviewByVariantID(tx, variantID); err != nil {
		tx.Rollback()
		return err
	}
	if err := DeleteCartByVariantID(tx, variantID); err != nil {
		tx.Rollback()
		return err
	}
	if err := DeleteOrderDetailByVariantID(tx, variantID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("variant_id = ?", variantID).Delete(&models.Variant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}
func GetImageByVariantID(variantID int32) (string, error) {
	var variant models.Variant
	if err := database.Db.Select("image").First(&variant, variantID).Error; err != nil {
		return "", err
	}
	return variant.Image, nil
}

func DeleteVariantByProductID(tx *gorm.DB, productID int32) error {
	var variants []models.Variant
	if err := tx.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return err
	}

	for _, variant := range variants {
		if err := DeleteReviewByVariantID(tx, variant.VariantID); err != nil {
			return err
		}
		if err := DeleteCartByVariantID(tx, variant.VariantID); err != nil {
			return err
		}
		if err := DeleteOrderDetailByVariantID(tx, variant.VariantID); err != nil {
			return err
		}
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.Variant{}).Error; err != nil {
		return err
	}

	return nil
}
