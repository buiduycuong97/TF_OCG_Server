package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

type OptionResult struct {
	OptionValueID   int32
	OptionProductID int32
	Value           string
	ProductID       int32
	OptionType      string
}

func CreateOptionProduct(optionProduct *models.OptionProduct) (*models.OptionProduct, error) {
	existingOptionProduct := &models.OptionProduct{}
	database.Db.Raw("SELECT * FROM option_product WHERE option_type = ?", optionProduct.OptionType).Scan(existingOptionProduct)
	if existingOptionProduct.OptionProductID != 0 {
		return nil, errors.New("Option Product already exists")
	}
	err := database.Db.Create(optionProduct).Error
	if err != nil {
		return nil, err
	}
	return optionProduct, nil
}

func GetAllOptionProduct() ([]*OptionResult, error) {
	var result []*OptionResult
	err := database.Db.Table("option_products").
		Select("option_products.*, option_values.*").
		Joins("LEFT JOIN option_values ON option_products.option_product_id = option_values.option_product_id").
		Scan(&result).
		Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetOptionProductByProductId(productID int32) ([]*OptionResult, error) {
	var result []*OptionResult
	err := database.Db.Table("option_products").
		Select("option_products.*, option_values.*").
		Joins("LEFT JOIN option_values ON option_products.option_product_id = option_values.option_product_id").
		Where("product_id = ?", productID).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetListOptionProductByProductID(productID int32) ([]models.OptionProduct, error) {
	var optionProducts []models.OptionProduct
	if err := database.Db.Where("product_id = ?", productID).Find(&optionProducts).Error; err != nil {
		return nil, err
	}
	return optionProducts, nil
}
