package dbms

import (
	"errors"
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateOptionValue(optionValue *models.OptionValue) (*models.OptionValue, error) {
	existingOptionValue := &models.OptionValue{}
	database.Db.Raw("SELECT * FROM option_value WHERE value = ?", optionValue.Value).Scan(existingOptionValue)
	if existingOptionValue.OptionValueID != 0 {
		return nil, errors.New("Option Value already exists")
	}
	err := database.Db.Create(optionValue).Error
	if err != nil {
		return nil, err
	}
	return optionValue, nil
}

func GetAllOptionValue() ([]*models.OptionValue, error) {
	optionValues := []*models.OptionValue{}
	err := database.Db.Find(&optionValues).Error
	if err != nil {
		return nil, err
	}
	return optionValues, nil
}

func GetOptionValueByOptionProductId(optionProductId int32) ([]*models.OptionValue, error) {
	optionValues := []*models.OptionValue{}
	err := database.Db.Where("option_product_id = ?", optionProductId).Find(&optionValues).Error
	if err != nil {
		return nil, err
	}
	return optionValues, nil
}

func GetOptionValueById(optionValue *models.OptionValue, optionValueID int32) error {
	result := database.Db.First(optionValue, optionValueID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteOptionValuesByOptionProductID(tx *gorm.DB, optionProductID int32) error {
	if err := tx.Where("option_product_id = ?", optionProductID).Delete(&models.OptionValue{}).Error; err != nil {
		return err
	}

	return nil
}

func DeleteOptionValuesByOptionProduct(optionProductId int32, optionValue string) error {
	if err := database.Db.Where("option_product_id = ? AND value = ?", optionProductId, optionValue).Delete(&models.OptionValue{}).Error; err != nil {
		return err
	}
	return nil
}
