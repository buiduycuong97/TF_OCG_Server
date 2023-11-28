package option_value_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func GetOptionValueByOptionProductId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	optionProductId, err := strconv.Atoi(vars["id"])

	optionValueList, err := dbms.GetOptionValueByOptionProductId(int32(optionProductId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.JSON(w, http.StatusOK, optionValueList)
}
