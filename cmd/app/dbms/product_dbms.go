package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateProduct(product *models.Product) (*models.Product, error) {
	existingProduct := &models.Product{}
	database.Db.Raw("SELECT * FROM products WHERE handle = ?", product.Handle).Scan(existingProduct)
	if existingProduct.Handle == product.Handle {
		return nil, errors.New("Product handle already exist")
	}
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	err := database.Db.Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetProductById(product *models.Product, id int32) (err error) {
	err = database.Db.Where("product_id = ?", id).First(product).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateProduct(updatedProduct *models.Product, id int32) error {
	database.Db.Model(updatedProduct).Where("product_id = ?", id).Updates(updatedProduct)
	return nil
}

func DeleteProduct(product *models.Product, id int32) error {
	database.Db.Where("product_id = ?", id).Updates(product)
	return nil
}

func GetListProduct(page int32, pageSize int32) ([]*models.Product, int64, error) {
	offset := (int64(page) - 1) * int64(pageSize)
	products := []*models.Product{}
	var totalCount int64

	if err := database.Db.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := database.Db.Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func GetListProductByCategoryId(categoryID int, page int32, pageSize int32) ([]*models.Product, int64, error) {
	offset := (page - 1) * pageSize
	products := []*models.Product{}
	var totalCount int64

	if err := database.Db.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := database.Db.Where("category_id = ?", categoryID).Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func SearchProduct(searchText string, categoryIDs []int, from string, to string, page int32, pageSize int32) ([]*models.Product, error) {
	offset := (page - 1) * pageSize
	products := []*models.Product{}

	query := database.Db

	if searchText != "" {
		query = query.Where("title LIKE ?", "%"+searchText+"%")
	}

	if len(categoryIDs) > 0 {
		query = query.Where("category_id IN (?)", categoryIDs)
	}

	if from != "" && to != "" {
		query = query.Where("price BETWEEN ? AND ?", from, to)
	}

	query = query.Offset(int(offset)).Limit(int(pageSize))

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
