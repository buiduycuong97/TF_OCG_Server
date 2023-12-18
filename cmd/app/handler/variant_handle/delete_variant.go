package variant_handle

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/product_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	variantID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var variant models.Variant
	err = dbms.GetVariantById(&variant, int32(variantID))

	var product models.Product
	product, err = dbms.GetProductByID(variant.ProductID)

	err = dbms.DeleteVariant(int32(variantID))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.DeleteProductFromCache(product_handle.RedisClient, product.Handle)
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Variant deleted successfully"})
}
