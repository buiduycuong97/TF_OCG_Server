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
