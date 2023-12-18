package variant_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	variantID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = dbms.DeleteVariant(int32(variantID))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Variant deleted successfully"})
}
