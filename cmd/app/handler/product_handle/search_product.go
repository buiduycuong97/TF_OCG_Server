package product_handle

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	categoryIDsStr := r.URL.Query().Get("categoryIDs")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	priceFrom := r.URL.Query().Get("priceFrom")
	priceTo := r.URL.Query().Get("priceTo")
	typeSort := r.URL.Query().Get("typeSort")
	fieldSort := r.URL.Query().Get("fieldSort")

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
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	products, err := dbms.SearchProduct(searchText, categoryIDs, priceFrom, priceTo, int32(page), int32(pageSize))

	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	SortProducts(products, typeSort, fieldSort)

	res.JSON(w, http.StatusOK, products)
}

func SortProducts(products []*models.Product, typeSort, fieldSort string) {
	switch fieldSort {
	case "title":
		if typeSort == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Title < products[j].Title
			})
		} else if typeSort == "desc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Title > products[j].Title
			})
		}
	case "price":
		if typeSort == "asc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price < products[j].Price
			})
		} else if typeSort == "desc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price > products[j].Price
			})
		}
	}
}
