package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

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

func GetAllOptionProduct() ([]*models.OptionProduct, error) {
	optionProducts := []*models.OptionProduct{}
	err := database.Db.Find(&optionProducts).Error
	if err != nil {
		return nil, err
	}
	return optionProducts, nil
}

func GetOptionProductByProductId(productID int32) ([]*models.OptionProduct, error) {
	optionProducts := []*models.OptionProduct{}
	err := database.Db.Where("product_id = ?", productID).Find(&optionProducts).Error
	if err != nil {
		return nil, err
	}
	return optionProducts, nil
}
