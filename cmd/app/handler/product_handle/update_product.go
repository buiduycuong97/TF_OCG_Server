package product_handle

import (
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"log"
	"net/http"
	"os"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/product_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var product models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if product.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Title is required"))
		return
	}

	productExist, err := dbms.GetProductByID(product.ProductID)

	product.Handle = productExist.Handle
	product.CategoryID = productExist.CategoryID
	product.CreatedAt = productExist.CreatedAt
	product.Image = productExist.Image

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

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("Error creating Elasticsearch client: %s", err)
		return
	}

	err = dbms.UpdateProduct(&product, esClient)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	//productData, err := convertProductToString(product)
	//if err != nil {
	//	log.Println("Lỗi chuyển đổi dữ liệu sản phẩm thành JSON: ", err)
	//	res.ERROR(w, http.StatusInternalServerError, err)
	//	return
	//}
	//
	//err = utils.SetProductToCache(RedisClient, product.Handle, productData)
	//if err != nil {
	//	log.Println("Lưu sản phẩm vào cache thất bại: ", err)
	//}

	//products, err := dbms.GetListProduct()
	//if err != nil {
	//	res.ERROR(w, http.StatusInternalServerError, err)
	//	return
	//}
	//cacheKey := "list_products"
	//err = utils.SetListProductsToCache(RedisClient, cacheKey, products)
	//if err != nil {
	//	log.Println("Lưu danh sách sản phẩm vào cache thất bại: ", err)
	//}

	productResponse := response.ProductResponseUpdate{
		ProductID:   product.ProductID,
		Handle:      product.Handle,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		UpdatedAt:   product.UpdatedAt,
	}

	err = utils.DeleteProductFromCache(RedisClient, product.Handle)
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	err = utils.DeleteListProductsFromCache(RedisClient, "list_products")
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, productResponse)
}
