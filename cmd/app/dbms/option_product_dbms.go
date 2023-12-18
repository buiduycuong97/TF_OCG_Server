package dbms

import (
	"errors"
	"gorm.io/gorm"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

type OptionResult struct {
	OptionValueID   int32  `json:"optionValueId"`
	OptionProductID int32  `json:"optionProductId"`
	Value           string `json:"value"`
	ProductID       int32  `json:"productId"`
	OptionType      string `json:"optionType"`
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
func DeleteOptionProductByProductID(tx *gorm.DB, productID int32) error {
	optionProductIDs, err := getOptionProductIDsByProductID(tx, productID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, optionProductID := range optionProductIDs {
		if err := DeleteOptionValuesByOptionProductID(tx, optionProductID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("product_id = ?", productID).Delete(&models.OptionProduct{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func getOptionProductIDsByProductID(tx *gorm.DB, productID int32) ([]int32, error) {
	var optionProductIDs []int32
	if err := tx.Model(&models.OptionProduct{}).Where("product_id = ?", productID).
		Pluck("option_product_id", &optionProductIDs).Error; err != nil {
		return nil, err
	}
	return optionProductIDs, nil
}

func GetOptionProductByOptionProductId(optionProductID int32) (*models.OptionProduct, error) {
	var optionProduct models.OptionProduct
	err := database.Db.Where("option_product_id = ?", optionProductID).First(&optionProduct).Error
	if err != nil {
		return nil, err
	}
	return &optionProduct, nil
}
