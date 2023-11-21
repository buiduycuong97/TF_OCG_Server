package product_handle

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	categoryIDsStr := r.URL.Query().Get("categoryIDs")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	priceFrom := r.URL.Query().Get("priceFrom")
	priceTo := r.URL.Query().Get("priceTo")

	// Kiểm tra xem searchText có được cung cấp không
	if searchText == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("searchText is required"))
		return
	}

	var categoryIDs []int
	if categoryIDsStr != "" {
		categoryIDStrArr := strings.Split(categoryIDsStr, ",")
		for _, idStr := range categoryIDStrArr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				res.ERROR(w, http.StatusBadRequest, err)
				return
			}
			categoryIDs = append(categoryIDs, id)
		}
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	products, err := dbms.SearchProduct(searchText, categoryIDs, priceFrom, priceTo, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, products)
}
