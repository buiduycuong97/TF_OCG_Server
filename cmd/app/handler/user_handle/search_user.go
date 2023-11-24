package user_handle

import (
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	users, err := dbms.SearchUser(searchText, int32(page), int32(pageSize))

	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, users)
}
