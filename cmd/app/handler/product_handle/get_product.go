package product_handle

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	option_product_response "tf_ocg/cmd/app/dto/option_product/response"
	option_value_response "tf_ocg/cmd/app/dto/option_value/response"
	"tf_ocg/cmd/app/dto/product_dto/response"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

var redisClient *redis.Client // Declare the Redis client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func convertProductToString(product models.Product) (string, error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)

	return jsonString, nil
}

func convertProductHandleToString(product response.ProductWithOptionResponse) (string, error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)

	return jsonString, nil
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	cachedData, err := utils.GetProductFromCache(redisClient, productID)
	if err == nil {
		res.JSON(w, http.StatusOK, cachedData)
		return
	}

	var product models.Product

	parsedProductID, err := strconv.ParseInt(productID, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = dbms.GetProductById(&product, int32(parsedProductID))
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

	err = utils.SetProductToCache(redisClient, productID, productData)
	if err != nil {
		log.Println("Lưu sản phẩm vào cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, product)
}

func GetProductByHandle(w http.ResponseWriter, r *http.Request) {
	handle := r.URL.Query().Get("handle")

	cachedData, err := utils.GetProductHandleFromCache(redisClient, handle)
	if err == nil {
		res.JSON(w, http.StatusOK, cachedData)
		return
	}

	var product models.Product
	err = dbms.GetProductByHandle(&product, handle)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	optionProducts, err := dbms.GetListOptionProductByProductID(product.ProductID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var optionProductsResponse []option_product_response.OptionProductResponse
	for _, optionSet := range optionProducts {
		var optionValuesResponse []option_value_response.OptionValueResponse
		optionValues, err := dbms.GetOptionValueByOptionProductId(optionSet.OptionProductID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		for _, optionValue := range optionValues {
			optionValuesResponse = append(optionValuesResponse, option_value_response.OptionValueResponse{
				OptionValueID:   optionValue.OptionValueID,
				OptionProductID: optionValue.OptionProductID,
				Value:           optionValue.Value,
			})
		}

		optionProductsResponse = append(optionProductsResponse, option_product_response.OptionProductResponse{
			OptionProductID: optionSet.OptionProductID,
			ProductID:       optionSet.ProductID,
			OptionType:      optionSet.OptionType,
			OptionValues:    optionValuesResponse,
		})
	}

	result := response.ProductWithOptionResponse{
		Product:        product,
		OptionProducts: optionProductsResponse,
	}

	productData, err := convertProductHandleToString(result)
	if err != nil {
		log.Println("Lỗi chuyển đổi dữ liệu sản phẩm thành JSON: ", err)
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.SetProductToCache(redisClient, handle, productData)
	if err != nil {
		log.Println("Lưu sản phẩm vào cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, result)
}

func GetListProducts(w http.ResponseWriter, r *http.Request) {
	cacheKey := "list_products"
	cachedData, err := utils.GetListProductsFromCache(redisClient, cacheKey)
	if err == nil {
		res.JSON(w, http.StatusOK, cachedData)
		return
	}

	products, err := dbms.GetListProduct()
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.SetListProductsToCache(redisClient, cacheKey, products)
	if err != nil {
		log.Println("Lưu danh sách sản phẩm vào cache thất bại: ", err)
	}

	res.JSON(w, http.StatusOK, products)
}

func GetListProductByCategoryId(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	categoryIDStr := r.URL.Query().Get("categoryId")

	categoryId, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid categoryId", http.StatusBadRequest)
		return
	}
	if categoryIDStr == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("categoryID is required"))
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

	products, totalCount, err := dbms.GetListProductByCategoryId(categoryId, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	result := map[string]interface{}{
		"products":   products,
		"totalPages": utils_handle.CalculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, result)
}
