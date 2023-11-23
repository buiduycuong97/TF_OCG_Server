package product_handle

import (
	"net/http"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	categoryIDsStr := r.URL.Query().Get("categoryIDs")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	priceFrom := r.URL.Query().Get("priceFrom")
	priceTo := r.URL.Query().Get("priceTo")

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

	if searchText == "" && categoryIDsStr == "" && pageStr == "" && pageSizeStr == "" && priceFrom == "" && priceTo == "" {
		allProducts, totalCount, err := dbms.GetListProduct(int32(page), int32(pageSize))
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		response := map[string]interface{}{
			"products":   allProducts,
			"totalPages": utils_handle.CalculateTotalPages(totalCount, int32(pageSize)),
		}
		res.JSON(w, http.StatusOK, response)
		return
	}

	products, err := dbms.SearchProduct(searchText, categoryIDs, priceFrom, priceTo, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, products)
}
