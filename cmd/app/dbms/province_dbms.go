package dbms

import (
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func GetProvinceByID(provinceID int32) (*models.Province, error) {
	var province models.Province
	err := database_manager.Db.Where("province_id = ?", provinceID).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func GetAllProvince(province *[]models.Province) (err error) {
	err = database_manager.Db.Order("province_name asc").Find(province).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProvinceByNameFromDB(provinceName string) (*models.Province, error) {
	var province models.Province
	err := database_manager.Db.Where("province_name = ?", provinceName).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func GetProvinceByIDFromDB(provinceID int32) (*models.Province, error) {
	var province models.Province
	err := database_manager.Db.Where("province_id = ?", provinceID).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}
