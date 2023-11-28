package option_product_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CreateOptionProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body"))
		return
	}
	var optionProduct models.OptionProduct
	err = json.Unmarshal(body, &optionProduct)
	if optionProduct.ProductID == 0 || optionProduct.OptionType == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ProductID or OptionType is empty"))
		return
	}

	result, err := dbms.CreateOptionProduct(&optionProduct)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusCreated, result)
}
