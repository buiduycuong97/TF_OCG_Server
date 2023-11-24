package category_handle

import (
	"errors"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func SearchCategories(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Kiểm tra xem searchText có được cung cấp không
	if searchText == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("searchText is required"))
		return
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

	categories, err := dbms.SearchCategory(searchText, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, categories)
}
