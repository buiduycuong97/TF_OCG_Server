package product_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	pid32 := int32(pid)

	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var product models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = dbms.UpdateProduct(&product, pid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, product)
}
