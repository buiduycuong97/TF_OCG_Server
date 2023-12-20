package product_handle

import (
	"cloud.google.com/go/storage"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/product_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func getFilePath() (string, error) {
	absPath, err := filepath.Abs("double2c-firebase-adminsdk-64fpu-cb7acf1b93.json")
	if err != nil {
		return "Fail to get file path", err
	} else {
		return absPath, nil
	}
}

const (
	//credentialsPath       = "./double2c-firebase-adminsdk-64fpu-cb7acf1b93.json"
	firebaseStorageBucket = "double2c.appspot.com"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	firebaseFilePath, error := getFilePath()
	if error != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, error)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Limit your maxMultipartMemory
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer file.Close()

	// Khởi tạo Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(firebaseFilePath))
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer client.Close()

	// Tạo đường dẫn trong bucket để lưu file
	objectPath := "images/" + handler.Filename

	// Tạo đối tượng để ghi dữ liệu vào Firebase Storage
	wc := client.Bucket(firebaseStorageBucket).Object(objectPath).NewWriter(ctx)
	if wc == nil {
		// Xử lý lỗi khi không lấy được đối tượng
		fmt.Printf("Error obtaining object writer: %v\n", err)
		return
	}

	// Copy dữ liệu từ file nhận được từ client vào đối tượng ghi (writer)
	if _, err := io.Copy(wc, file); err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = wc.Close()
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Lấy thông tin chi tiết về tệp đã lưu trữ
	attrs := wc.Attrs()
	if attrs == nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// attrs.Path sẽ chứa đường dẫn lưu trữ của tệp
	//storagePath := attrs.Name

	// Tạo URL để truy cập file đã lưu trong Firebase Storage
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", attrs.Bucket, attrs.Name)

	product.Title = r.PostFormValue("title")
	product.Description = r.PostFormValue("description")
	product.Price, _ = strconv.ParseFloat(r.PostFormValue("price"), 64)
	product.CategoryID, _ = strconv.Atoi(r.PostFormValue("categoryID"))
	product.Image = imageURL

	if product.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Title is required"))
		return
	}
	product.Handle = slug.Make(product.Title)

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

	_, err = dbms.CreateProduct(&product, esClient)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	createProductRes := response.Product{
		ProductId:   product.ProductID,
		Handle:      product.Handle,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		Image:       product.Image,
	}

	err = utils.DeleteListProductsFromCache(RedisClient, "list_products")
	if err != nil {
		log.Println("Xóa sản phẩm trong cache thất bại: ", err)
	}

	res.JSON(w, http.StatusCreated, createProductRes)
}
