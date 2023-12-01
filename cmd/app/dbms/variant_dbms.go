package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateVariant(variant *models.Variant) (*models.Variant, error) {
	existingVariant := &models.Variant{}
	database.Db.Raw("SELECT * FROM variants WHERE product_id = ? AND option_product1 = ? AND option_product2 = ?", variant.ProductID, variant.OptionProduct1, variant.OptionProduct2).Scan(existingVariant)
	if existingVariant.VariantID > 0 {
		return nil, errors.New("Variant with the same product, optionProduct1, and optionProduct2 already exists")
	}
	err := database.Db.Create(variant).Error
	if err != nil {
		return nil, err
	}
	return variant, nil
}
