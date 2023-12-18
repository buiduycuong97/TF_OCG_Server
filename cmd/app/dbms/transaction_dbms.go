package dbms

import (
	"gorm.io/gorm"
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	transaction.CreatedAt = time.Now()
	err := database_manager.Db.Create(transaction).Error
	return transaction, err
}

func DeleteTransactionByOrderId(tx *gorm.DB, orderID int32) error {
	// Delete Transactions associated with the OrderID
	if err := tx.Where("order_id = ?", orderID).Delete(&models.Transaction{}).Error; err != nil {
		return err
	}

	return nil
}
