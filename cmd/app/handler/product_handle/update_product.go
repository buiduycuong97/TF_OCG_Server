package product_handle

import (
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"log"
	"net/http"
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
		log.Printf("Error creating Elasticsearch client: %s", err)
		return
	}

	err = dbms.UpdateProduct(&product, esClient)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	productData, err := convertProductToString(product)
	if err != nil {
		log.Println("Lỗi chuyển đổi dữ liệu sản phẩm thành JSON: ", err)
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.SetProductToCache(redisClient, product.Handle, productData)
	if err != nil {
		log.Println("Lưu sản phẩm vào cache thất bại: ", err)
	}

	products, err := dbms.GetListProduct()
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	cacheKey := "list_products"
	err = utils.SetListProductsToCache(redisClient, cacheKey, products)
	if err != nil {
		log.Println("Lưu danh sách sản phẩm vào cache thất bại: ", err)
	}

	productResponse := response.ProductResponseUpdate{
		ProductID:   product.ProductID,
		Handle:      product.Handle,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		UpdatedAt:   product.UpdatedAt,
	}

	res.JSON(w, http.StatusOK, productResponse)
}
