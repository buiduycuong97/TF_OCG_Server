package option_value_handle

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/product_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func CreateOptionValue(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body"))
		return
	}
	var optionValue models.OptionValue
	err = json.Unmarshal(body, &optionValue)
	if optionValue.OptionProductID == 0 || optionValue.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("OptionProductID or Value is empty"))
		return
	}

	result, err := dbms.CreateOptionValue(&optionValue)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var optionProdcut *models.OptionProduct
	optionProdcut, err = dbms.GetOptionProductByOptionProductId(optionValue.OptionProductID)

	var product models.Product
	product, err = dbms.GetProductByID(optionProdcut.ProductID)

	err = utils.DeleteProductFromCache(product_handle.RedisClient, product.Handle)
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	res.JSON(w, http.StatusCreated, result)
}
