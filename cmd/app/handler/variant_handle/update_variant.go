package variant_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func UpdateVariantQuantityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	variantID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = UpdateVariantCountInStock(int32(variantID), int32(quantity))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Variant quantity updated successfully"})
}

func UpdateVariantCountInStock(variantID int32, newQuantity int32) error {
	variant := &models.Variant{}
	err := dbms.GetVariantById(variant, variantID)
	if err != nil {
		return errors.New("Failed to get variant")
	}

	if variant.CountInStock < newQuantity {
		return errors.New("Not enough quantity remaining")
	}

	variant.CountInStock -= newQuantity

	err = dbms.UpdateVariant(variant, variantID)
	if err != nil {
		return errors.New("Failed to update variant quantity")
	}

	return nil
}

func UpdateVariantQuantityWithIncrease(variantID int32, quantityToIncrease int32) error {
	variant := &models.Variant{}
	err := dbms.GetVariantById(variant, variantID)
	if err != nil {
		return errors.New("Failed to get variant")
	}

	variant.CountInStock += quantityToIncrease

	err = dbms.UpdateVariant(variant, variantID)
	if err != nil {
		return errors.New("Failed to update variant quantity")
	}

	return nil
}
