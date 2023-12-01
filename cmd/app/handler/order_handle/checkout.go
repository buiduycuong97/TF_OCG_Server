// Trong package order_handle hoặc một file tương tự
package order_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"tf_ocg/cmd/app/dbms"
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

	// Khai báo một biến để giải mã JSON
	var requestData map[string]interface{}

	// Giải mã dữ liệu JSON vào biến requestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, fmt.Errorf("Failed to unmarshal JSON: %v", err))
		return
	}

	// Kiểm tra xem requestData có chứa các trường cần thiết không
	shippingAddress, ok := requestData["shippingAddress"].(string)
	if !ok || shippingAddress == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Shipping address is required"))
		return
	}

	provinceID, ok := requestData["provinceId"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Province ID is required"))
		return
	}

	// Chuyển đổi giá trị provinceID thành int32
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

	//Thêm trường totalQuantity vào requestData
	totalQuantity, ok := requestData["totalQuantity"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Total quantity is required"))
		return
	}

	convertedTotalQuantity := int32(totalQuantity)

	// Thêm trường totalPrice vào requestData
	totalPrice, ok := requestData["totalPrice"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Total price is required"))
		return
	}

	// Thêm trường grandTotal vào requestData
	grandTotal, ok := requestData["grandTotal"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Grand total is required"))
		return
	}

	// Thêm trường discountAmount vào requestData
	discountAmount, ok := requestData["discountAmount"].(float64)
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("Discount amount is required"))
		return
	}

	newOrder := models.Order{
		UserID:          userID,
		OrderDate:       time.Now(),
		ShippingAddress: shippingAddress,
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

	tx.Commit()

	// Gọi hàm cập nhật số lượng trong kho của biến thể sau khi checkout thành công
	for _, cartItem := range cartItems {
		err := variant_handle.UpdateVariantCountInStock(cartItem.VariantID, cartItem.Quantity)
		if err != nil {
			// Xử lý lỗi nếu cần
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
