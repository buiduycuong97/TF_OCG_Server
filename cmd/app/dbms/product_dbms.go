package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateProduct(product *models.Product) (*models.Product, error) {
	existingProduct := &models.Product{}
	database.Db.Raw("SELECT * FROM products WHERE title = ?", product.Title).Scan(existingProduct)
	if existingProduct.Handle == product.Handle {
		return nil, errors.New("Product title already exist")
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
	err = database.Db.Where("product_id = ?", id).Find(product).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProductByHandle(product *models.Product, handle string) (err error) {
	err = database.Db.Where("handle = ?", handle).Find(product).Error
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
	if err := deleteReviewsByProductID(id); err != nil {
		return err
	}
	if err := deleteCartItemByProductID(id); err != nil {
		return err
	}
	if err := deleteOrderDetailByProductID(id); err != nil {
		return err
	}
	return database.Db.Where("product_id = ?", id).Delete(product).Error
}

func deleteCartItemByProductID(productID int32) error {
	return database.Db.Where("product_id = ?", productID).Delete(&models.Cart{}).Error
}

func deleteReviewsByProductID(productID int32) error {
	return database.Db.Where("product_id = ?", productID).Delete(&models.Review{}).Error
}

func deleteOrderDetailByProductID(productID int32) error {
	return database.Db.Where("product_id = ?", productID).Delete(&models.OrderDetail{}).Error
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

func SearchProduct(searchText string, categories []string, from string, to string, page int32, pageSize int32, typeSort string, fieldSort string) ([]*models.Product, int64, error) {
	offset := (page - 1) * pageSize
	products := []*models.Product{}
	var totalCount int64

	query := database.Db.Model(&models.Product{})

	if searchText != "" {
		query = query.Where("title LIKE ?", "%"+searchText+"%")
	}

	if len(categories) > 0 {
		query = query.Joins("JOIN categories ON products.category_id = categories.category_id").
			Where("categories.handle IN (?)", categories)
	}

	if from != "" && to != "" {
		query = query.Where("price BETWEEN ? AND ?", from, to)
	}

	if err := query.Find(&[]*models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if fieldSort != "" && typeSort != "" {
		query = query.Order(fieldSort + " " + typeSort)
	}

	query = query.Offset(int(offset)).Limit(int(pageSize))

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func UpdateProductQuantity(productID int32, newQuantity int32) error {
	product := &models.Product{}
	err := GetProductById(product, productID)
	if err != nil {
		return errors.New("Failed to get product")
	}

	if product.QuantityRemaining < newQuantity {
		return errors.New("Not enough quantity remaining")
	}

	product.QuantityRemaining -= newQuantity

	err = UpdateProduct(product, productID)
	if err != nil {
		return errors.New("Failed to update product quantity")
	}

	return nil
}

func UpdateProductQuantityWithIncrease(productID int32, quantityToIncrease int32) error {
	product := &models.Product{}
	err := GetProductById(product, productID)
	if err != nil {
		return errors.New("Failed to get product")
	}

	product.QuantityRemaining += quantityToIncrease

	err = UpdateProduct(product, productID)
	if err != nil {
		return errors.New("Failed to update product quantity")
	}

	return nil
}

func GetProductByID(productID int32) (models.Product, error) {
	var product models.Product
	err := database.Db.First(&product, productID).Error
	return product, err
}
