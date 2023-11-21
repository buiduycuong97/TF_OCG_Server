package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateCategory(category *models.Category) (*models.Category, error) {
	existingCategory := &models.Category{}
	database.Db.Raw("SELECT * FROM category WHERE name = ?", category.Handle).Scan(existingCategory)
	if existingCategory.Handle == category.Handle {
		return nil, errors.New("Category name already exist")
	}
	err := database.Db.Create(category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func GetCategoryById(category *models.Category, id int32) (err error) {
	err = database.Db.Where("category_id = ?", id).First(category).Error
	if err != nil {
		return err
	}
	return nil
}

func GetListCategory(page int32, pageSize int32) ([]*models.Category, int64, error) {
	offset := (int64(page) - 1) * int64(pageSize)
	categories := []*models.Category{}
	var totalCount int64

	if err := database.Db.Model(&models.Category{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := database.Db.Offset(int(offset)).Limit(int(pageSize)).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}

func UpdateCategory(updatedCategory *models.Category, id int32) error {
	database.Db.Model(updatedCategory).Where("category_id = ?", id).Updates(updatedCategory)
	return nil
}

func DeleteCategory(category *models.Category, id int32) error {
	database.Db.Where("category_id = ?", id).Updates(category)
	return nil
}

func SearchCategory(searchText string, page int32, pageSize int32) ([]*models.Category, error) {
	offset := (page - 1) * pageSize
	categories := []*models.Category{}

	query := database.Db

	if searchText != "" {
		query = query.Where("title LIKE ?", "%"+searchText+"%")
	}

	query = query.Offset(int(offset)).Limit(int(pageSize))

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
