package option_value_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/product_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func DeleteOptionValue(w http.ResponseWriter, r *http.Request) {
	var opValue models.OptionValue
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &opValue)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	vars := mux.Vars(r)
	optionProductId, err := strconv.ParseUint(vars["id"], 10, 32)

	var optionProdcut *models.OptionProduct
	optionProdcut, err = dbms.GetOptionProductByOptionProductId(int32(optionProductId))

	var product models.Product
	product, err = dbms.GetProductByID(optionProdcut.ProductID)

	err = dbms.DeleteOptionValuesByOptionProduct(int32(optionProductId), opValue.Value)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
	}

	err = utils.DeleteProductFromCache(product_handle.RedisClient, product.Handle)
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}
	res.JSON(w, http.StatusOK, "Delete option value success")

}
