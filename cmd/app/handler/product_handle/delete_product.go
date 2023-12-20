package product_handle

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
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

	if err := godotenv.Load(); err != nil {
		return
	}

	Addresses := os.Getenv("ES_ADDRESS")
	Username := os.Getenv("ES_USERNAME")
	Password := os.Getenv("ES_PASSWORD")
	// Cấu hình Elasticsearch
	esCfg := elasticsearch.Config{
		Addresses: []string{Addresses},
		Username:  Username,
		Password:  Password,
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

	var product models.Product
	err = dbms.GetProductById(&product, pid32)

	var handle = product.Handle

	// Xóa sản phẩm và cập nhật Elasticsearch
	err = dbms.DeleteProductDB(esClient, pid32)
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

	err = utils.DeleteProductFromCache(RedisClient, handle)
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	err = utils.DeleteListProductsFromCache(RedisClient, "list_products")
	if err != nil {
		log.Println("Xóa list sản phẩm trong cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, data)
}
