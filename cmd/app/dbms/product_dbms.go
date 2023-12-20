package dbms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"io"
	"log"
	"strconv"
	"strings"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/proto/models"
	"time"
)

func CreateProduct(product *models.Product, esClient *elasticsearch.Client) (*models.Product, error) {
	existingProduct := &models.Product{}
	database.Db.Raw("SELECT * FROM products WHERE title = ?", product.Title).Scan(existingProduct)
	if existingProduct.Handle == product.Handle {
		return nil, errors.New("Product title already exists")
	}

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// Thêm sản phẩm vào cơ sở dữ liệu
	tx := database.Db.Begin()
	if err := tx.Create(product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	// Kiểm tra và chỉ mục sản phẩm trong Elasticsearch
	if err := IndexProductES(esClient, product); err != nil {
		// Nếu có lỗi khi chỉ mục sản phẩm, rollback thêm sản phẩm trong cơ sở dữ liệu
		tx := database.Db.Begin()
		if err := tx.Where("product_id = ?", product.ProductID).Delete(&models.Product{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()

		return nil, err
	}

	return product, nil
}

func IndexProductES(esClient *elasticsearch.Client, product *models.Product) error {
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
		return err
	}

	dataJSON, err := json.Marshal(indexData)
	if err != nil {
		log.Printf("Error marshalling product data to JSON: %s", err)
		return err
	}

	bulkRequestBody := []string{string(indexJSON), string(dataJSON)}

	bulkRequestString := strings.Join(bulkRequestBody, "\n") + "\n"

	resp, err := esClient.Bulk(strings.NewReader(bulkRequestString))
	if err != nil {
		log.Printf("Error sending bulk request to Elasticsearch: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Elasticsearch response body: %s", err)
		return err
	}

	if resp.IsError() {
		log.Printf("Elasticsearch responded with error: %s", resp.Status())
		log.Printf("Response body: %s", body)
		return errors.New("Elasticsearch response error")
	}

	log.Printf("Indexing completed for product with ID %d", product.ProductID)
	return nil
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

func GetProductById(product *models.Product, id int32) (err error) {
	err = database.Db.Where("product_id = ?", id).Find(product).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProductByHandle(product *models.Product, handle string) (err error) {
	err = database.Db.Where("handle = ?", handle).Find(product).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateProduct(product *models.Product, esClient *elasticsearch.Client) error {
	existingProduct := &models.Product{}
	err := database.Db.Where("product_id = ?", product.ProductID).Find(existingProduct).Error
	if err != nil {
		return err
	}

	if existingProduct.Handle != product.Handle {
		return errors.New("Product does not exist")
	}

	product.Handle = slug.Make(product.Title)
	if err := godotenv.Load(); err != nil {
		return err
	}

	tx := database.Db.Begin()
	if err := tx.Model(product).Where("product_id = ?", product.ProductID).Updates(product).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// Cập nhật sản phẩm trong Elasticsearch
	if err := UpdateProductES(esClient, product); err != nil {
		// Nếu có lỗi khi cập nhật Elasticsearch, rollback cập nhật trong cơ sở dữ liệu quan hệ
		tx := database.Db.Begin()
		if err := tx.Model(product).Where("product_id = ?", product.ProductID).Updates(existingProduct).Error; err != nil {
			log.Printf("Error rolling back database update: %s", err)
		}
		tx.Commit()

		return err
	}

	return nil
}

func UpdateProductES(esClient *elasticsearch.Client, product *models.Product) error {
	log.Printf("Updating product in Elasticsearch. Product ID: %d", product.ProductID)

	updateData := map[string]interface{}{
		"doc": map[string]interface{}{
			"handle":      product.Handle,
			"title":       product.Title,
			"description": convertToValidJSON(cleanDescription(product.Description)),
			"price":       product.Price,
			"categoryID":  product.CategoryID,
			"image":       product.Image,
			"created_at":  product.CreatedAt,
			"updated_at":  product.UpdatedAt,
		},
	}

	updateJSON, err := json.Marshal(updateData)
	if err != nil {
		log.Printf("Error marshalling update request to JSON: %s", err)
		return err
	}

	log.Printf("Update JSON: %s", updateJSON)

	resp, err := esClient.Update(
		"products",
		strconv.Itoa(int(product.ProductID)),
		strings.NewReader(string(updateJSON)),
		esClient.Update.WithContext(context.Background()),
		esClient.Update.WithRefresh("true"),
	)

	if err != nil {
		log.Printf("Error sending update request to Elasticsearch: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Elasticsearch response body: %s", err)
		return err
	}

	log.Printf("Elasticsearch response status: %s", resp.Status())
	log.Printf("Response body: %s", body)

	if resp.IsError() {
		log.Printf("Elasticsearch responded with error: %s", resp.Status())
		return errors.New("Elasticsearch response error")
	}

	log.Printf("Update completed for product with ID %d in Elasticsearch", product.ProductID)
	return nil
}

func DeleteProductDB(esClient *elasticsearch.Client, id int32) error {
	err := DeleteProductES(esClient, id)
	if err != nil {
		return err
	}

	err = DeleteProductFromDB(id)
	if err != nil {
		ReAddProductToES(esClient, id)
		return err
	}

	return nil
}

func ReAddProductToES(esClient *elasticsearch.Client, productID int32) error {
	product, err := GetProductByID(productID)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin sản phẩm từ cơ sở dữ liệu: %s", err)
		return err
	}

	err = IndexProductES(esClient, &product)
	if err != nil {
		log.Printf("Lỗi khi chỉ mục lại sản phẩm trong Elasticsearch: %s", err)
		return err
	}

	log.Printf("Thêm lại sản phẩm có ID %d vào Elasticsearch thành công", productID)
	return nil
}

func DeleteProductFromDB(id int32) error {
	tx := database.Db.Begin()

	if err := DeleteOptionProductByProductID(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := DeleteVariantByProductID(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("product_id = ?", id).Delete(&models.Product{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func DeleteProductES(esClient *elasticsearch.Client, id int32) error {
	resp, err := esClient.Delete(
		"products",
		strconv.Itoa(int(id)),
		esClient.Delete.WithContext(context.Background()),
		esClient.Delete.WithRefresh("true"),
	)

	if err != nil {
		log.Printf("Error sending delete request to Elasticsearch: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Elasticsearch response body: %s", err)
		return err
	}

	if resp.IsError() {
		log.Printf("Elasticsearch responded with error: %s", resp.Status())
		log.Printf("Response body: %s", body)
		return errors.New("Elasticsearch response error")
	}

	log.Printf("Delete completed for product with ID %d in Elasticsearch", id)
	return nil
}

func GetListProduct() ([]*models.Product, error) {
	products := []*models.Product{}

	err := database.Db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetListProductByCategoryId(categoryID int, page int32, pageSize int32) ([]*models.Product, int64, error) {
	offset := (page - 1) * pageSize
	products := []*models.Product{}
	var totalCount int64

	if err := database.Db.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := database.Db.Where("category_id = ?", categoryID).Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func SearchProduct(searchText string, categories []string, from string, to string, page int32, pageSize int32, typeSort string, fieldSort string) ([]*models.Product, int64, error) {
	offset := (page - 1) * pageSize
	products := []*models.Product{}
	var totalCount int64

	query := database.Db.Model(&models.Product{})

	if searchText != "" {
		query = query.Where("title LIKE ?", "%"+searchText+"%")
	}

	if len(categories) > 0 {
		query = query.Joins("JOIN categories ON products.category_id = categories.category_id").
			Where("categories.handle IN (?)", categories)
	}

	if from != "" && to != "" {
		query = query.Where("price BETWEEN ? AND ?", from, to)
	}

	if err := query.Find(&[]*models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if fieldSort != "" && typeSort != "" {
		query = query.Order(fieldSort + " " + typeSort)
	}

	query = query.Offset(int(offset)).Limit(int(pageSize))

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func GetProductByID(productID int32) (models.Product, error) {
	var product models.Product
	err := database.Db.First(&product, productID).Error
	return product, err
}

func SearchProductES(
	esClient *elasticsearch.Client,
	searchText string,
	categoryIDs []int32,
	priceFrom, priceTo string,
	page, pageSize int32,
	typeSort, fieldSort string,
) (products []models.Product, totalItems int, err error) {
	query := buildElasticsearchQuery(searchText, categoryIDs, priceFrom, priceTo, page, pageSize, typeSort, fieldSort)
	res, err := esClient.Search(
		esClient.Search.WithIndex("products"),
		esClient.Search.WithBody(strings.NewReader(query)),
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		log.Printf("Error performing Elasticsearch search: %s", err)
		return nil, 0, err
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %s", closeErr)
		}
	}()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, 0, err
	}

	if res.IsError() {
		log.Printf("Response body: %v", response)

		return nil, 0, errors.New("Elasticsearch response error")
	}

	hits, _ := response["hits"].(map[string]interface{})
	totalItems = int(hits["total"].(map[string]interface{})["value"].(float64))

	hitsArray, _ := hits["hits"].([]interface{})
	for _, hit := range hitsArray {
		source, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})

		createdAt, _ := convertToTime(source["created_at"].(string))
		updatedAt, _ := convertToTime(source["updated_at"].(string))

		product := models.Product{
			ProductID:   int32(source["productId"].(float64)),
			Handle:      source["handle"].(string),
			Title:       source["title"].(string),
			Description: source["description"].(string),
			Price:       source["price"].(float64),
			CategoryID:  int(source["categoryID"].(float64)),
			Image:       source["image"].(string),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		products = append(products, product)
	}

	return products, totalItems, nil
}

func buildElasticsearchQuery(
	searchText string,
	categoryIDs []int32,
	priceFrom, priceTo string,
	page, pageSize int32,
	typeSort, fieldSort string,
) string {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
		"size": pageSize,
		"from": int(page-1) * int(pageSize),
	}

	if searchText != "" {
		// Thêm điều kiện tìm kiếm theo title
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"match_phrase_prefix": map[string]interface{}{
					"title": searchText,
				},
			},
		)
	}

	if len(categoryIDs) > 0 {
		categoryFilter := map[string]interface{}{
			"terms": map[string]interface{}{
				"categoryID": categoryIDs,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			categoryFilter,
		)
	}

	if priceFrom != "" && priceTo != "" {
		priceRangeFilter := map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": priceFrom,
					"lte": priceTo,
				},
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			priceRangeFilter,
		)
	}

	// Thêm điều kiện sắp xếp
	if fieldSort == "" {
		sort := map[string]interface{}{
			"updatedAt": map[string]interface{}{
				"order": "desc", // Sắp xếp theo updatedAt theo thứ tự giảm dần
			},
		}
		query["sort"] = []interface{}{sort}
	} else if fieldSort != "" && typeSort != "" {
		if fieldSort == "title" {
			fieldSort = "title.keyword"
		}
		sort := map[string]interface{}{
			fieldSort: map[string]interface{}{
				"order": typeSort,
			},
		}
		query["sort"] = []interface{}{sort}
	}

	queryJSON, _ := json.Marshal(query)
	return string(queryJSON)
}

func convertToTime(timestamp string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.999Z"
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
