package product_handle

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/utils"

	"github.com/elastic/go-elasticsearch/v8"
)

func SetElasticsearchClient(client *elasticsearch.Client) {
	utils.EsClient = client
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
	categoryStr := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	priceFrom := r.URL.Query().Get("priceFrom")
	priceTo := r.URL.Query().Get("priceTo")
	typeSort := r.URL.Query().Get("typeSort")
	fieldSort := r.URL.Query().Get("fieldSort")

	var handleStrs []string
	var categoryIDs []int32

	if categoryStr != "" {
		handleStrs = strings.Split(categoryStr, ",")

		allCategories, err := dbms.GetAllCategories()
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		categoryIDMap := make(map[string]int32)
		for _, cat := range allCategories {
			categoryIDMap[cat.Handle] = cat.CategoryID
		}

		for _, handle := range handleStrs {
			if categoryID, found := categoryIDMap[handle]; found {
				categoryIDs = append(categoryIDs, categoryID)
			}
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

	esCfg := elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "Ksckb67MQwA-frPDAA7+",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	products, totalItems, err := dbms.SearchProductES(esClient, searchText, categoryIDs, priceFrom, priceTo, int32(page), int32(pageSize), typeSort, fieldSort)

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
