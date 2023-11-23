package dbms

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func AddToCart(userID, productID, quantity int32) (*models.Cart, error) {
	tx := database_manager.Db.Begin()

	existingCart := &models.Cart{}
	err := tx.Where("user_id = ? AND product_id = ?", userID, productID).First(existingCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Record not found, creating new item")
			cart, err := createNewCartItem(userID, productID, quantity, tx)
			if err != nil {
				tx.Rollback()
				fmt.Println("Error creating new item:", err)
				return nil, err
			}
			tx.Commit()
			return cart, nil
		}
		tx.Rollback()
		return nil, err
	}

	existingCart.Quantity += quantity
	err = tx.Model(existingCart).Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", existingCart.Quantity).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Error updating existing cart:", err)
		return nil, err
	}

	tx.Commit()
	return existingCart, nil
}

func createNewCartItem(userID, productID, quantity int32, tx *gorm.DB) (*models.Cart, error) {
	product := &models.Product{}
	err := GetProductById(product, productID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	totalPrice := product.Price * float64(quantity)

	newCartItem := &models.Cart{
		UserID:     userID,
		ProductID:  productID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
	}

	err = tx.Create(newCartItem).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return newCartItem, nil
}

func GetCartByUserID(userID int32) ([]models.Cart, error) {
	var cartItems []models.Cart
	err := database_manager.Db.Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(cartItems); i++ {

	}
	return cartItems, nil
}

func UpdateCartItem(userID int32, productID, quantity int) error {

	cartItem, err := GetCartItem(userID, productID)
	if err != nil {
		return errors.New("Failed to get cart item")
	}

	quantityInt32 := int32(quantity)

	cartItem.Quantity = quantityInt32
	cartItem.TotalPrice = float64(quantityInt32)
	err = UpdateCart(cartItem, cartItem.CartID)
	if err != nil {
		return errors.New("Failed to update cart item")
	}

	return nil
}

func UpdateCart(cart *models.Cart, cartID int32) error {
	err := database_manager.Db.Model(&models.Cart{}).Where("cart_id = ?", cartID).Updates(cart).Error
	if err != nil {
		return errors.New("Failed to update cart in the database")
	}
	return nil
}

func GetCartItem(userID int32, productID int) (*models.Cart, error) {
	cartItem := &models.Cart{}
	err := database_manager.Db.Where("user_id = ? AND product_id = ?", userID, productID).First(cartItem).Error
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func RemoveCartItem(userID int32, productID int) error {
	cartItem, err := GetCartItem(userID, int(int32(productID)))
	if err != nil {
		return errors.New("Failed to get cart item")
	}

	err = DeleteCart(cartItem, cartItem.CartID)
	if err != nil {
		return errors.New("Failed to remove item from cart in the database")
	}

	return nil
}

func DeleteCart(cart *models.Cart, cartID int32) error {
	err := database_manager.Db.Model(&models.Cart{}).Where("cart_id = ?", cartID).Delete(cart).Error
	if err != nil {
		return errors.New("Failed to delete cart item from the database")
	}
	return nil
}
