package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateCategory(categoryCreate *models.Categories) (*models.Categories, error) {
	existingCategory := &models.Categories{}
	database.Db.Raw("SELECT * FROM categories WHERE name = ?", categoryCreate.Name).Scan(existingCategory)
	if existingCategory.Handle == categoryCreate.Handle {
		return nil, errors.New("Category name already exist")
	}
	err := database.Db.Create(categoryCreate).Error
	if err != nil {
		return nil, err
	}
	return categoryCreate, nil
}

func GetCategoryById(category *models.Categories, id int32) (err error) {
	err = database.Db.Where("category_id = ?", id).First(category).Error
	if err != nil {
		return err
	}
	return nil
}

func GetListCategory(page int32, pageSize int32) ([]*models.Categories, int64, error) {
	offset := (int64(page) - 1) * int64(pageSize)
	categories := []*models.Categories{}
	var totalCount int64

	if err := database.Db.Model(&models.Categories{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := database.Db.Offset(int(offset)).Limit(int(pageSize)).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}

func GetAllCategories() ([]*models.Categories, error) {
	categories := []*models.Categories{}

	if err := database.Db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func UpdateCategory(updatedCategory *models.Categories, id int32) error {
	database.Db.Model(updatedCategory).Where("category_id = ?", id).Updates(updatedCategory)
	return nil
}

func DeleteCategory(category *models.Categories, id int32) error {
	var products []models.Product
	database.Db.Where("category_id = ?", id).Find(&products)
	for _, product := range products {
		if err := database.Db.Delete(&product).Error; err != nil {
			return err
		}
	}
	if err := database.Db.Where("category_id = ?", id).Delete(category).Error; err != nil {
		return err
	}

	return nil
}

func SearchCategory(searchText string, page int32, pageSize int32) ([]*models.Categories, error) {
	offset := (page - 1) * pageSize
	categories := []*models.Categories{}

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
