package option_value_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
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

	res.JSON(w, http.StatusCreated, result)
}
