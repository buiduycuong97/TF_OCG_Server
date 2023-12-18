package dbms

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
)

func AddToCart(userID, variantID, quantity int32) (*models.Cart, error) {
	tx := database_manager.Db.Begin()

	existingCart := &models.Cart{}
	err := tx.Where("user_id = ? AND variant_id = ?", userID, variantID).First(existingCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Record not found, creating new item")
			cart, err := createNewCartItem(userID, variantID, quantity, tx)
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
	err = tx.Model(existingCart).Where("user_id = ? AND variant_id = ?", userID, variantID).
		Update("quantity", existingCart.Quantity).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Error updating existing cart:", err)
		return nil, err
	}

	tx.Commit()
	return existingCart, nil
}

func createNewCartItem(userID, variantID, quantity int32, tx *gorm.DB) (*models.Cart, error) {
	variant := &models.Variant{}
	err := GetVariantById(variant, variantID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	totalPrice := float64(variant.Price) * float64(quantity)

	newCartItem := &models.Cart{
		UserID:     userID,
		VariantID:  variantID,
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

func UpdateCartItem(userID, variantID, quantity int) error {
	cartItem, err := GetCartItem(int32(userID), variantID)
	if err != nil {
		return errors.New("Failed to get cart item")
	}

	variant, err := GetVariant(variantID)
	if err != nil {
		return errors.New("Failed to get variant details")
	}

	cartItem.Quantity = int32(quantity)
	cartItem.TotalPrice = float64(quantity) * float64(variant.Price)

	if err := UpdateCart(cartItem, cartItem.CartID); err != nil {
		return errors.New("Failed to update cart item")
	}

	return nil
}

func GetVariant(variantID int) (*models.Variant, error) {
	variant := &models.Variant{}
	err := GetVariantById(variant, int32(variantID))
	return variant, err
}

func UpdateCart(cart *models.Cart, cartID int32) error {
	err := database_manager.Db.Model(&models.Cart{}).Where("cart_id = ?", cartID).Updates(cart).Error
	if err != nil {
		return errors.New("Failed to update cart in the database")
	}
	return nil
}

func GetCartItem(userID int32, variantID int) (*models.Cart, error) {
	cartItem := &models.Cart{}
	err := database_manager.Db.Where("user_id = ? AND variant_id = ?", userID, variantID).First(cartItem).Error
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func RemoveCartItem(userID int32, variantID int) error {
	cartItem, err := GetCartItem(userID, variantID)
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

func DeleteCartByVariantID(tx *gorm.DB, variantID int32) error {
	if err := tx.Where("variant_id = ?", variantID).Delete(&models.Cart{}).Error; err != nil {
		return err
	}

	return nil
}
