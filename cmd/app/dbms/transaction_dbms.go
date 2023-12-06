package dbms

import (
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	transaction.CreatedAt = time.Now()
	err := database_manager.Db.Create(transaction).Error
	return transaction, err
}
