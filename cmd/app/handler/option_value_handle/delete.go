package option_value_handle

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

	err = dbms.DeleteOptionValuesByOptionProduct(int32(optionProductId), opValue.Value)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
	}
	res.JSON(w, http.StatusOK, "Delete option value success")
}
