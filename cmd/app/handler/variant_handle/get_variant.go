package variant_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func GetVariantIdByOption(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID    int32 `json:"productID"`
		OptionValue1 int32 `json:"optionValue1"`
		OptionValue2 int32 `json:"optionValue2"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}

	variantID, err := dbms.GetVariantIdByOption(request.ProductID, request.OptionValue1, request.OptionValue2)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response_api.JSON(w, http.StatusOK, map[string]int32{"variantId": variantID})
}

func GetListVariantByOrderId(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}

	orderDetails, err := dbms.GetOrderDetailsByOrderID(int32(orderID))
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var variantDetails []map[string]interface{}

	for _, orderDetail := range orderDetails {
		variantDetail, err := GetVariantDetails(orderDetail.VariantID, orderDetail.OrderDetailID, orderDetail.IsReview)
		if err != nil {
			response_api.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		variantDetails = append(variantDetails, variantDetail)
	}

	responselist := map[string]interface{}{
		"variants": variantDetails,
	}

	response_api.JSON(w, http.StatusOK, responselist)
}

func GetVariantDetails(variantID int32, orderDetailID int32, isReview bool) (map[string]interface{}, error) {
	var variant models.Variant
	var product models.Product
	var option1 models.OptionValue
	var option2 models.OptionValue

	// Lấy thông tin biến thể
	err := dbms.GetVariantById(&variant, variantID)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin sản phẩm từ ProductID trong biến thể
	err = dbms.GetProductById(&product, variant.ProductID)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin OptionValue1 (nếu có)
	if variant.OptionValue1 != 0 {
		err = dbms.GetOptionValueById(&option1, variant.OptionValue1)
		if err != nil {
			return nil, err
		}
	}

	// Lấy thông tin OptionValue2 (nếu có)
	if variant.OptionValue2 != 0 {
		err = dbms.GetOptionValueById(&option2, variant.OptionValue2)
		if err != nil {
			return nil, err
		}
	}

	// Tạo response theo định dạng mong muốn
	response := map[string]interface{}{
		"variantId": variant.VariantID,
		"productId": map[string]interface{}{
			"productId": product.ProductID,
			"title":     product.Title,
		},
		"title":        variant.Title,
		"price":        variant.Price,
		"comparePrice": variant.ComparePrice,
		"countInStock": variant.CountInStock,
		"image":        variant.Image,
	}

	// Thêm OptionValue1 vào response (nếu có)
	if variant.OptionValue1 != 0 {
		response["optionValue1"] = map[string]interface{}{
			"optionValueId":   option1.OptionValueID,
			"optionProductId": option1.OptionProductID,
			"value":           option1.Value,
		}
	}

	// Thêm OptionValue2 vào response (nếu có)
	if variant.OptionValue2 != 0 {
		response["optionValue2"] = map[string]interface{}{
			"optionValueId":   option2.OptionValueID,
			"optionProductId": option2.OptionProductID,
			"value":           option2.Value,
		}
	}

	response["orderDetailId"] = orderDetailID
	response["isReview"] = isReview

	return response, nil
}

func GetVariantById(w http.ResponseWriter, r *http.Request) {
	variantIDStr := mux.Vars(r)["id"]
	variantID, err := strconv.ParseInt(variantIDStr, 10, 32)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var variant models.Variant
	err = dbms.GetVariantById(&variant, int32(variantID))
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response_api.JSON(w, http.StatusOK, variant)
}
