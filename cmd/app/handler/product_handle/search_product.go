package product_handle

import (
	"net/http"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	categoryStr := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	priceFrom := r.URL.Query().Get("priceFrom")
	priceTo := r.URL.Query().Get("priceTo")
	typeSort := r.URL.Query().Get("typeSort")
	fieldSort := r.URL.Query().Get("fieldSort")

	var categories []string
	if categoryStr != "" {
		categories = strings.Split(categoryStr, ",")
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	products, totalItems, err := dbms.SearchProduct(searchText, categories, priceFrom, priceTo, int32(page), int32(pageSize), typeSort, fieldSort)

	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"products":   products,
		"totalItems": totalItems,
	}

	res.JSON(w, http.StatusOK, response)
}
