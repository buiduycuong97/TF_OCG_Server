package category_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	cid32 := int32(cid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var category models.Categories
	err = dbms.DeleteCategory(&category, cid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	data := struct {
		Message string `json:"message"`
	}{
		"Category deleted successfully",
	}
	res.JSON(w, http.StatusOK, data)
}
