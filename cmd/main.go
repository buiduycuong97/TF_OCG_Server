package main

import (
	// ... (import statements)

	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"tf_ocg/cmd/app/handler/discount_handle"
	"tf_ocg/cmd/app/handler/order_handle"
	"tf_ocg/cmd/app/handler/product_handle"
	"tf_ocg/cmd/app/handler/utils_handle"
	"tf_ocg/cmd/app/router"
	"tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

type Server struct {
	Db                  *gorm.DB
	Router              *mux.Router
	RedisClient         *redis.Client
	ElasticsearchClient *elasticsearch.Client
}

func Init() {
	var server Server
	server.Db = database_manager.InitDb()
	server.Router = mux.NewRouter()
	router.InitializeRoutes(server.Router)
	order_handle.ScheduleOrderStatusUpdate()
	discount_handle.ScheduleDiscountCodeGeneration()

	config, err := utils.LoadConfig("./config.json")
	if err != nil {
		log.Fatal("Không thể đọc cấu hình: ", err)
	}

	server.RedisClient = utils.NewRedisClient(config)
	product_handle.SetRedisClient(server.RedisClient)

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

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	product_handle.SetElasticsearchClient(esClient)
	server.ElasticsearchClient = esClient

	//createIndexMapping(Addresses, Username, Password)

	var products []models.Product
	server.Db.Find(&products)

	log.Printf("Number of products retrieved from MySQL: %d", len(products))

	var bulkRequestBody []string
	for _, product := range products {
		indexData := map[string]interface{}{
			"productId":   product.ProductID,
			"handle":      product.Handle,
			"title":       product.Title,
			"description": convertToValidJSON(cleanDescription(product.Description)),
			"price":       product.Price,
			"categoryID":  product.CategoryID,
			"image":       product.Image,
			"created_at":  product.CreatedAt,
			"updated_at":  product.UpdatedAt,
		}

		indexJSON, err := json.Marshal(map[string]interface{}{"index": map[string]interface{}{"_index": "products", "_id": product.ProductID}})
		if err != nil {
			log.Printf("Error marshalling index request to JSON: %s", err)
			continue
		}

		dataJSON, err := json.Marshal(indexData)
		if err != nil {
			log.Printf("Error marshalling product data to JSON: %s", err)
			continue
		}

		bulkRequestBody = append(bulkRequestBody, string(indexJSON), string(dataJSON))
	}

	bulkRequestString := strings.Join(bulkRequestBody, "\n") + "\n"

	log.Printf("Starting bulk indexing for products")

	resp, err := esClient.Bulk(strings.NewReader(bulkRequestString))
	if err != nil {
		log.Printf("Error sending bulk request to Elasticsearch: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Elasticsearch response body: %s", err)
		return
	}

	if resp.IsError() {
		log.Printf("Elasticsearch responded with error: %s", resp.Status())
		log.Printf("Response body: %s", body)
		return
	}

	log.Printf("Bulk indexing completed for products")

	go utils_handle.HandleRabbitMQMessages()
	server.Run(":8000")

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}

	port := os.Getenv("SERVER_PORT")
	server.Run(port)
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, cors.AllowAll().Handler(server.Router)))
}

func convertToValidJSON(description string) string {
	type DescriptionJSON struct {
		Description string `json:"description"`
	}

	descJSON := DescriptionJSON{Description: description}
	jsonBytes, err := json.Marshal(descJSON)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return ""
	}

	return string(jsonBytes)
}

func cleanDescription(description string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r == '"' || r == '\'' || r == ',' {
			return -1
		}
		return r
	}, description)

	return cleaned
}

//func createIndexMapping(address, username, password string) {
//	// Cấu hình kết nối Elasticsearch
//	cfg := elasticsearch.Config{
//		Addresses: []string{address},
//		Username:  username,
//		Password:  password,
//		Transport: &http.Transport{
//			TLSClientConfig: &tls.Config{
//				InsecureSkipVerify: true,
//			},
//		},
//	}
//
//	// Khởi tạo client Elasticsearch
//	esClient, err := elasticsearch.NewClient(cfg)
//	if err != nil {
//		log.Fatalf("Error creating Elasticsearch client: %s", err)
//	}
//
//	// Khởi tạo request PUT mapping
//	indexRequest := map[string]interface{}{
//		"mappings": map[string]interface{}{
//			"properties": map[string]interface{}{
//				"productId":   map[string]interface{}{"type": "integer"},
//				"handle":      map[string]interface{}{"type": "keyword"},
//				"title":       map[string]interface{}{"type": "keyword"},
//				"description": map[string]interface{}{"type": "text"},
//				"price":       map[string]interface{}{"type": "float"},
//				"categoryID":  map[string]interface{}{"type": "integer"},
//				"image":       map[string]interface{}{"type": "keyword"},
//				"createdAt":   map[string]interface{}{"type": "date"},
//				"updatedAt":   map[string]interface{}{"type": "date"},
//			},
//		},
//	}
//
//	req := esClient.Indices.Create("products").Body(strings.NewReader(serializeToJson(indexRequest)))
//
//	// Thực hiện request PUT mapping
//	res, err := req.Do(context.Background())
//	if err != nil {
//		log.Fatalf("Error creating index mapping: %s", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		log.Fatalf("Error creating index mapping: %s", res.Status())
//	}
//
//	log.Println("Index mapping created successfully")
//}

func serializeToJson(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func main() {
	Init()
}
