package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateVariant(variant *models.Variant) (*models.Variant, error) {
	existingVariant := &models.Variant{}
	database.Db.Raw("SELECT * FROM variants WHERE product_id = ? AND option_value1 = ? AND option_product2 = ?", variant.ProductID, variant.OptionValue1, variant.OptionValue2).Scan(existingVariant)
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
	return database.Db.Where("variant_id = ?", variantID).First(variant).Error
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
