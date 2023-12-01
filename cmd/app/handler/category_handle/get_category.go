package category_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	cid32 := int32(cid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var category models.Categories
	err = dbms.GetCategoryById(&category, cid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, category)
}

func GetListCategories(w http.ResponseWriter, r *http.Request) {
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

	categories, totalCount, err := dbms.GetListCategory(int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"categories": categories,
		"totalPages": calculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, response)
}

func calculateTotalPages(totalCount int64, pageSize int32) int32 {
	totalPages := int32(totalCount) / pageSize
	if int32(totalCount)%pageSize != 0 {
		totalPages++
	}
	return totalPages
}
