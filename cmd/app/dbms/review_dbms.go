package dbms

import (
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func CreateReview(review *models.Review) error {
	return database_manager.Db.Create(review).Error
}

func GetReviewsByProductID(productID int32) ([]models.Review, error) {
	var reviews []models.Review
	err := database_manager.Db.Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}
