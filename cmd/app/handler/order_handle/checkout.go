// Trong package order_handle hoặc một file tương tự
package order_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/discount_handle"
	"tf_ocg/cmd/app/handler/utils_handle"
	"tf_ocg/cmd/app/handler/variant_handle"
	database "tf_ocg/pkg/database_manager"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"time"
)

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	user, err := dbms.GetUserByID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, fmt.Errorf("Failed to read request body: %v", err))
		return
	}

	var requestData map[string]interface{}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, fmt.Errorf("Failed to unmarshal JSON: %v", err))
		return
	}

	shippingAddress, ok := requestData["shippingAddress"].(string)
	if !ok || shippingAddress == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Shipping address is required"))
		return
	}

	phoneOrder, ok := requestData["phoneOrder"].(string)
	if !ok || phoneOrder == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Phone order is required"))
		return
	}

	provinceID, ok := requestData["provinceId"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Province ID is required"))
		return
	}

	convertedProvinceID := int32(provinceID)

	cartItems, err := dbms.GetCartByUserID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(cartItems) == 0 {
		res.ERROR(w, http.StatusBadRequest, errors.New("Cart is empty, cannot proceed with checkout"))
		return
	}

	totalQuantity, ok := requestData["totalQuantity"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Total quantity is required"))
		return
	}

	convertedTotalQuantity := int32(totalQuantity)

	totalPrice, ok := requestData["totalPrice"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Total price is required"))
		return
	}

	grandTotal, ok := requestData["grandTotal"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Grand total is required"))
		return
	}

	discountAmount, ok := requestData["discountAmount"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Discount amount is required"))
		return
	}

	discountCode, ok := requestData["discountCode"].(string)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Discount amount is required"))
		return
	}
	if discountCode != "" {
		var discount models.Discount
		if err := dbms.GetDiscountByDiscountCodeAndUserID(&discount, discountCode, int(userID)); err != nil {
			res.ERROR(w, http.StatusNotFound, errors.New("Invalid discount code"))
			return
		}

		// Giảm số lượng discount đi 1
		discount.AvailableQuantity--
		if err := dbms.UpdateDiscount(&discount, discount.DiscountID); err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if discount.AvailableQuantity == 0 || time.Now().After(discount.EndDate) {
			if err := dbms.DeleteDiscount(&discount, discount.DiscountID); err != nil {
				res.ERROR(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	newOrder := models.Order{
		UserID:          userID,
		OrderDate:       time.Now(),
		ShippingAddress: shippingAddress,
		PhoneOrder:      phoneOrder,
		Status:          models.Pending,
		ProvinceID:      convertedProvinceID,
		TotalQuantity:   convertedTotalQuantity,
		TotalPrice:      totalPrice,
		GrandTotal:      grandTotal,
		DiscountAmount:  discountAmount,
	}

	tx := database.Db.Begin()

	createdOrder, err := dbms.CreateOrder(tx, &newOrder)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Created Order ID:", createdOrder.OrderID)

	province, err := dbms.GetProvinceByID(newOrder.ProvinceID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	for _, cartItem := range cartItems {
		orderDetail := &models.OrderDetail{
			OrderID:   createdOrder.OrderID,
			VariantID: cartItem.VariantID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.TotalPrice,
		}

		err = dbms.CreateOrderDetail(tx, orderDetail)
		if err != nil {
			tx.Rollback()
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	newOrder.TotalPrice -= newOrder.DiscountAmount

	err = dbms.ClearCart(tx, userID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to clear cart: %v", err))
		return
	}

	err = dbms.UpdateOrderTotalValues(tx, createdOrder.OrderID, newOrder.TotalQuantity, newOrder.TotalPrice)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to update order values: %v", err))
		return
	}

	updateUserStatsAndCheckDiscount(tx, user, createdOrder)

	tx.Commit()

	for _, cartItem := range cartItems {
		err := variant_handle.UpdateVariantCountInStock(cartItem.VariantID, cartItem.Quantity)
		if err != nil {
		}
	}

	responseData := map[string]interface{}{
		"message":          "Order created successfully",
		"user_name":        user.UserName,
		"phone_number":     user.PhoneNumber,
		"province_name":    province.ProvinceName,
		"shipping_fee":     province.ShippingFee,
		"order_id":         createdOrder.OrderID,
		"order_date":       createdOrder.OrderDate,
		"shipping_address": createdOrder.ShippingAddress,
		"total_quantity":   newOrder.TotalQuantity,
		"total_price":      newOrder.TotalPrice,
		"discount_amount":  newOrder.DiscountAmount,
	}

	res.JSON(w, http.StatusCreated, responseData)

}

func updateUserStatsAndCheckDiscount(tx *gorm.DB, user *models.User, order *models.Order) {
	user.OrderCount++
	user.TotalSpent += int32(order.GrandTotal)

	err := UpdateLevelAndCheckDiscount(tx, user)
	if err != nil {
		return
	}
}

func UpdateLevelAndCheckDiscount(tx *gorm.DB, user *models.User) error {
	if user.OrderCount >= 3 && user.TotalSpent >= 1000000 && user.CurrentLevel == models.Bronze {
		user.CurrentLevel = models.Silver
		user.NextLevel = models.Gold
		discount, err := discount_handle.CreateAutomaticDiscountForUpgrade(user)
		if err != nil {
			return err
		}
		err = SendOrderStatusUpdateEmail(user.Email, string(user.CurrentLevel), discount.DiscountCode)
		if err != nil {
			return err
		}
	}

	if user.OrderCount >= 20 && user.TotalSpent >= 5000000 && user.CurrentLevel == models.Silver {
		user.CurrentLevel = models.Gold
		user.NextLevel = models.Diamond
		discount, err := discount_handle.CreateAutomaticDiscountForUpgrade(user)
		if err != nil {
			return err
		}
		err = SendOrderStatusUpdateEmail(user.Email, string(user.CurrentLevel), discount.DiscountCode)
		if err != nil {
			return err
		}
	}

	if user.OrderCount >= 75 && user.TotalSpent >= 15000000 && user.CurrentLevel == models.Gold {
		user.CurrentLevel = models.Diamond
		user.NextLevel = ""
		discount, err := discount_handle.CreateAutomaticDiscountForUpgrade(user)
		if err != nil {
			return err
		}
		err = SendOrderStatusUpdateEmail(user.Email, string(user.CurrentLevel), discount.DiscountCode)
		if err != nil {
			return err
		}
	}

	if err := dbms.UpdateUserLevel(tx, user, user.UserID); err != nil {
		return err
	}

	return nil
}

func SendOrderStatusUpdateEmail(email, currentLevel string, discount string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailhost := os.Getenv("EMAIL_HOST")
	subject := "Order Status Update"
	body := fmt.Sprintf("Congratulations on achieving %s membership, Double 2C will send you a discount code [%s]!", currentLevel, discount)

	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(emailhost, 587, emailAddress, emailPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
