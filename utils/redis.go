package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io/ioutil"
	"log"
	"tf_ocg/cmd/app/dto/product_dto/response"
	"tf_ocg/proto/models"
)

type Config struct {
	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
}

func LoadConfig(filePath string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func NewRedisClient(config Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("Không thể kết nối đến Redis server: ", err)
	}

	return client
}

func GetProductHandleFromCache(client *redis.Client, handle string) (*response.ProductWithOptionResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", handle)

	cachedData, err := client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("Cache miss for key:", cacheKey)
		return nil, err
	} else if err != nil {
		fmt.Println("Error retriev9ing data from cache:", err)
		return nil, err
	}

	var productWithOption response.ProductWithOptionResponse
	err = json.Unmarshal([]byte(cachedData), &productWithOption)
	if err != nil {
		fmt.Println("Error during JSON unmarshaling:", err)
		fmt.Println("Cached Data:", cachedData)
		return nil, err
	}

	fmt.Println("Successfully unmarshaled data from cache")
	return &productWithOption, nil
}

func SetProductToCache(client *redis.Client, handle string, productData string) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", handle)

	err := client.Set(ctx, cacheKey, productData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteProductFromCache(client *redis.Client, handle string) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", handle)

	err := client.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetListProductsFromCache(client *redis.Client, listProductCacheKey string) ([]models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("listProductCacheKey:%s", listProductCacheKey)

	cachedData, err := client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	var products []models.Product
	err = json.Unmarshal([]byte(cachedData), &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func SetListProductsToCache(client *redis.Client, listProductCacheKey string, products []*models.Product) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("listProductCacheKey:%s", listProductCacheKey)

	jsonData, err := json.Marshal(products)
	if err != nil {
		return err
	}

	err = client.Set(ctx, cacheKey, jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteListProductsFromCache(client *redis.Client, listProductCacheKey string) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("listProductCacheKey:%s", listProductCacheKey)

	// Hàm Del để xóa cache của danh sách sản phẩm
	err := client.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}
