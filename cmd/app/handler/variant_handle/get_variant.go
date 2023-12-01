package variant_handle

import (
	"encoding/json"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/pkg/response_api"
)

func GetVariantIdByOption(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID    int32 `json:"productID"`
		OptionValue1 int32 `json:"optionValue1"`
		OptionValue2 int32 `json:"optionValue2"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}

	variantID, err := dbms.GetVariantIdByOption(request.ProductID, request.OptionValue1, request.OptionValue2)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response_api.JSON(w, http.StatusOK, map[string]int32{"variantId": variantID})
}
