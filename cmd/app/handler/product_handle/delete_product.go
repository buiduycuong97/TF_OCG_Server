package product_handle

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Lấy thông tin sản phẩm từ URL
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	pid32 := int32(pid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Cấu hình Elasticsearch
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

	// Tạo client Elasticsearch
	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("Lỗi khi tạo client Elasticsearch: %s", err)
		return
	}

	// Xóa sản phẩm và cập nhật Elasticsearch
	err = dbms.DeleteProduct(esClient, pid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Trả về thông báo JSON thành công
	data := struct {
		Message string `json:"message"`
	}{
		"Sản phẩm đã được xóa thành công",
	}
	res.JSON(w, http.StatusOK, data)
}
