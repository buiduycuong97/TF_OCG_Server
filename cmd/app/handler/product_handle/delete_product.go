package product_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	pid32 := int32(pid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var product models.Product
	err = dbms.DeleteProduct(&product, pid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	data := struct {
		Message string `json:"message"`
	}{
		"Product deleted successfully",
	}
	res.JSON(w, http.StatusOK, data)
}
