package dbms

import (
	"gorm.io/gorm"
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateReview(review *models.Review) error {
	review.CreatedAt = time.Now()
	return database_manager.Db.Create(review).Error
}

func DeleteReviewByVariantID(tx *gorm.DB, variantID int32) error {
	// Delete Reviews associated with the VariantID
	if err := tx.Where("variant_id = ?", variantID).Delete(&models.Review{}).Error; err != nil {
		return err
	}
	return nil
}

func GetReviewsByVariantID(variantID int32) ([]models.Review, error) {
	var reviews []models.Review
	err := database_manager.Db.Where("variant_id = ?", variantID).Find(&reviews).Error
	return reviews, err
}
